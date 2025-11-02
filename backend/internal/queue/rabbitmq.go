package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"wishlist-go/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type Message struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}

var Client *RabbitMQClient

// ConnectRabbitMQ устанавливает соединение с RabbitMQ
func ConnectRabbitMQ() error {
	cfg := config.Config.RabbitMQ

	// Формируем URL подключения
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open a channel: %w", err)
	}

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
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	Client = &RabbitMQClient{
		conn:    conn,
		channel: ch,
		queue:   q,
	}

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

// PublishMessage отправляет сообщение в очередь
func (c *RabbitMQClient) PublishMessage(msgType string, payload map[string]interface{}) error {
	if c == nil {
		return fmt.Errorf("RabbitMQ client is not initialized")
	}

	msg := Message{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.channel.PublishWithContext(
		ctx,
		"",           // exchange
		c.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Message sent: type=%s, payload=%v", msgType, payload)
	return nil
}

// Close закрывает соединение с RabbitMQ
func (c *RabbitMQClient) Close() error {
	if c == nil {
		return nil
	}

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	log.Println("RabbitMQ connection closed")
	return nil
}
