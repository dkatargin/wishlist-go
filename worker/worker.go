package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/yaml.v3"
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

var config Config

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Загружаем конфигурацию
	if err := loadConfig(*configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключаемся к RabbitMQ
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Создаём канал
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Объявляем очередь
	q, err := ch.QueueDeclare(
		"wishlist_tasks", // имя очереди
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Устанавливаем prefetch count для равномерного распределения
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	// Начинаем получать сообщения
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Worker started. Waiting for messages...")

	// Обработка graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Горутина для обработки сообщений
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

	// Ждём сигнала завершения
	<-sigChan
	log.Println("Shutting down gracefully...")
	cancel()

	// Даём время на завершение обработки текущих сообщений
	time.Sleep(2 * time.Second)
	log.Println("Worker stopped")
}

// loadConfig загружает конфигурацию из файла
func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

// connectRabbitMQ устанавливает соединение с RabbitMQ
func connectRabbitMQ() (*amqp.Connection, error) {
	cfg := config.RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return conn, nil
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
	case "wishlist_created":
		handleWishlistCreated(msg.Payload)
	case "wishitem_created":
		handleWishItemCreated(msg.Payload)
	case "account_updated":
		handleAccountUpdated(msg.Payload)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}

	// Отправляем acknowledgment
	if err := delivery.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}

// handleWishlistCreated обрабатывает создание wishlist
func handleWishlistCreated(payload map[string]interface{}) {
	log.Printf("Processing wishlist_created: %v", payload)
	// Здесь может быть логика отправки уведомлений, обновления кэша и т.д.
}

// handleWishItemCreated обрабатывает создание wishitem
func handleWishItemCreated(payload map[string]interface{}) {
	log.Printf("Processing wishitem_created: %v", payload)
	// Здесь может быть логика отправки уведомлений, обновления кэша и т.д.
}

// handleAccountUpdated обрабатывает обновление account
func handleAccountUpdated(payload map[string]interface{}) {
	log.Printf("Processing account_updated: %v", payload)
	// Здесь может быть логика синхронизации данных, отправки email и т.д.
}
