package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"wishlist-api/internal/config"
	"wishlist-api/internal/db"
	"wishlist-api/internal/db/models"
	"wishlist-api/internal/queue"

	"wishlist-worker/crawler"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config структура для конфигурации worker
type Config struct {
	RabbitMQ struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Vhost    string `yaml:"vhost"`
	} `yaml:"rabbitmq"`
}

// Message структура входящего сообщения
type Message struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Загружаем конфигурацию и подключаемся к базе данных

	err := config.LoadConfigFile(*configPath)
	if err != nil {
		log.Panicf("Failed to load config file %s: %v+", *configPath, err.Error())
	}

	err = db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Подключаемся к RabbitMQ
	err = queue.ConnectRabbitMQ()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queue.Client.Close()

	// Начинаем получать сообщения из очереди
	msgs, err := queue.Client.ConsumeMessages()
	if err != nil {
		log.Panicf("Failed to start consuming messages: %v", err)
	}

	log.Println("Worker started. Waiting for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Channel closed")
					return
				}
				handleMessage(msg)
			}
		}
	}()

	<-sigChan
	log.Println("Shutting down gracefully...")
	cancel()

	time.Sleep(2 * time.Second)
	log.Println("Worker stopped")
}

// handleMessage обрабатывает входящее сообщение
func handleMessage(delivery amqp.Delivery) {
	var msg Message
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		// Отправляем negative acknowledgment
		delivery.Nack(false, false)
		return
	}

	log.Printf("Received message: Type=%s, Payload=%v, Timestamp=%v",
		msg.Type, msg.Payload, msg.Timestamp)

	// Здесь добавьте логику обработки различных типов сообщений
	switch msg.Type {
	case "crawl_product":
		handleCrawlProduct(msg.Payload, delivery)
		return // Возвращаемся раньше, т.к. ack/nack обрабатывается внутри
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}

	// Отправляем acknowledgment
	if err := delivery.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}

// handleCrawlProduct обрабатывает запрос на краулинг конкретного товара по URL
func handleCrawlProduct(payload map[string]interface{}, delivery amqp.Delivery) {
	productURL, ok := payload["product_url"].(string)
	if !ok || productURL == "" {
		log.Printf("Invalid product_url in payload: %v", payload)
		_ = delivery.Nack(false, false)
		return
	}

	log.Printf("Crawling product: %s", productURL)

	client := crawler.NewYaMarketClient()
	productInfo, err := client.FetchProductByURL(productURL)
	if err != nil {
		log.Printf("Failed to crawl product %s: %v", productURL, err)
		_ = delivery.Nack(false, true) // Отправляем обратно в очередь для повторной попытки
		return
	}

	// Логируем результат
	resultPayload := map[string]interface{}{
		"request_id":  payload["request_id"],
		"product_url": productInfo.URL,
		"name":        productInfo.Title,
		"price":       productInfo.Price,
		"currency":    "RUB",
		"image_url":   productInfo.ImageURL,
		"crawled_at":  time.Now(),
	}

	jsonData, _ := json.MarshalIndent(resultPayload, "", "  ")
	log.Printf("Product crawled successfully:\n%s", string(jsonData))

	go func(payload map[string]interface{}) {
		var crawlRequestInfo models.WishItemDataRequest
		err = db.ORM.Model(models.WishItemDataRequest{}).Where("request_id = ?", payload["request_id"]).First(&crawlRequestInfo).Error
		if err != nil {
			fmt.Errorf("failed to find crawl request ID: %s", payload["request_id"])
			return
		}
		productPriceFloat, err := strconv.ParseFloat(productInfo.Price, 64)
		if err != nil {
			fmt.Errorf("failed to parse product price: %s", productInfo.Price)
			return
		}
		// Сохраняем найденный товар в список желаний
		err = db.ORM.Model(models.WishItem{}).Create(&models.WishItem{
			WishListCode:   crawlRequestInfo.WishListCode,
			OwnerID:        crawlRequestInfo.OwnerID,
			Name:           productInfo.Title,
			Priority:       1,
			Status:         "pending",
			MarketLink:     productInfo.URL,
			MarketPicture:  productInfo.ImageURL,
			MarketPrice:    productPriceFloat,
			MarketCurrency: "RUB", // TODO: определить валюту в зависимости от сайта
			MarketQuantity: 1,
		}).Error
		if err != nil {
			fmt.Errorf("failed to save crawled product to DB: %v", err)
			return
		}

		// Обновляем статус для запроса
		err = db.ORM.Model(&crawlRequestInfo).Where("id = ?", crawlRequestInfo.ID).Updates(models.WishItemDataRequest{
			Status:    "completed",
			UpdatedAt: time.Now().Unix(),
		}).Error
		if err != nil {
			fmt.Errorf("failed to update crawl request status: %v", err)
			return
		}
	}(payload)

	_ = delivery.Ack(false)
}
