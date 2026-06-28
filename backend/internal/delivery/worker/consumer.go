package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/infrastructure/crawler"
	"wishlist-go/internal/infrastructure/queue"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer слушает очередь задач и обрабатывает сообщения воркера.
type Consumer struct {
	mq       *queue.RabbitMQClient
	crawler  *crawler.YaMarketClient
	itemRepo domain.WishItemRepository
}

// NewConsumer собирает consumer с инъекцией зависимостей.
func NewConsumer(mq *queue.RabbitMQClient, c *crawler.YaMarketClient, itemRepo domain.WishItemRepository) *Consumer {
	return &Consumer{mq: mq, crawler: c, itemRepo: itemRepo}
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

// handleCrawlProduct тянет товар краулером и сохраняет WishItem.
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
		nackBounded(d)
		return
	}

	item, err := crawlPayloadToWishItem(payload, info)
	if err != nil {
		log.Printf("worker: битый payload crawl_product %v: %v", payload, err)
		_ = d.Nack(false, false) // не реквеуим — payload не починится
		return
	}

	if err := w.itemRepo.CreateWishItem(item); err != nil {
		log.Printf("worker: не удалось сохранить WishItem: %v", err)
		nackBounded(d)
		return
	}

	log.Printf("worker: товар сохранён %s -> %q", productURL, info.Title)
	_ = d.Ack(false)
}

// crawlPayloadToWishItem собирает доменный WishItem из payload задачи и результата краула.
func crawlPayloadToWishItem(payload map[string]interface{}, info *crawler.ProductInfo) (*domain.WishItem, error) {
	codeStr, _ := payload["wish_list_code"].(string)
	shareCode, err := uuid.Parse(codeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid wish_list_code %q: %w", codeStr, err)
	}
	// owner_id: JSON number → float64 → int64. Точно для id ≤ 2^53; Telegram гарантирует ≤52 бита.
	ownerF, ok := payload["owner_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing/invalid owner_id")
	}

	price, _ := parsePrice(info.Price) // цена не критична: 0 при неудаче парсинга
	name := info.Title
	pic := info.ImageURL
	marketURL := info.URL
	return &domain.WishItem{
		WishListCode:     shareCode,
		OwnerID:          int64(ownerF),
		Name:             &name,
		MarketURL:        marketURL,
		MarketPictureURL: &pic,
		MarketPrice:      &price,
		MarketCurrency:   "RUB",
	}, nil
}

// nackBounded реквеуит сообщение не более одного раза (по d.Redelivered), затем дропает —
// чтобы «ядовитое» сообщение (битый URL / удалённый список) не крутилось в очереди вечно.
// TODO: заменить на DLX с лимитом попыток и отложенным ретраем.
func nackBounded(d amqp.Delivery) {
	_ = d.Nack(false, !d.Redelivered)
}

// parsePrice извлекает число из строки цены краулера (например "12345 ₽").
func parsePrice(s string) (float64, error) {
	var b strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return 0, fmt.Errorf("no numeric price in %q", s)
	}
	return strconv.ParseFloat(b.String(), 64)
}
