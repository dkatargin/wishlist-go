package worker

import (
	"context"
	"encoding/json"
	"log"
	"wishlist-go/internal/infrastructure/crawler"
	"wishlist-go/internal/infrastructure/queue"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer слушает очередь задач и обрабатывает сообщения воркера.
type Consumer struct {
	mq      *queue.RabbitMQClient
	crawler *crawler.YaMarketClient
}

// NewConsumer собирает consumer с инъекцией зависимостей.
func NewConsumer(mq *queue.RabbitMQClient, c *crawler.YaMarketClient) *Consumer {
	return &Consumer{mq: mq, crawler: c}
}

// Run потребляет сообщения до отмены ctx или закрытия канала очереди.
func (w *Consumer) Run(ctx context.Context) error {
	msgs, err := w.mq.ConsumeMessages()
	if err != nil {
		return err
	}
	log.Println("worker: ожидание сообщений...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-msgs:
			if !ok {
				log.Println("worker: канал очереди закрыт")
				return nil
			}
			w.handle(d)
		}
	}
}

func (w *Consumer) handle(d amqp.Delivery) {
	var msg queue.Message
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		log.Printf("worker: не удалось разобрать сообщение: %v", err)
		_ = d.Nack(false, false)
		return
	}

	switch msg.Type {
	case "crawl_product":
		w.handleCrawlProduct(msg.Payload, d)
	default:
		log.Printf("worker: неизвестный тип сообщения %q", msg.Type)
		_ = d.Ack(false)
	}
}

// handleCrawlProduct тянет товар краулером. Персист в WishItem (producer-флоу +
// репозиторий) ещё не подключён; сейчас результат логируется.
func (w *Consumer) handleCrawlProduct(payload map[string]interface{}, d amqp.Delivery) {
	productURL, _ := payload["product_url"].(string)
	if productURL == "" {
		log.Printf("worker: crawl_product без product_url: %v", payload)
		_ = d.Nack(false, false)
		return
	}

	info, err := w.crawler.FetchProductByURL(productURL)
	if err != nil {
		log.Printf("worker: краул %s не удался: %v", productURL, err)
		_ = d.Nack(false, true) // вернуть в очередь для повторной попытки
		return
	}

	log.Printf("worker: товар получен %s -> title=%q price=%q", productURL, info.Title, info.Price)
	_ = d.Ack(false)
}
