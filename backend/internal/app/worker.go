package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"wishlist-go/internal/delivery/worker"
	"wishlist-go/internal/infrastructure/config"
	"wishlist-go/internal/infrastructure/crawler"
	"wishlist-go/internal/infrastructure/queue"
)

// WorkerApp — приложение фонового воркера (consumer очереди задач).
type WorkerApp struct {
	consumer *worker.Consumer
	mq       *queue.RabbitMQClient
}

// NewWorkerApp собирает зависимости воркера.
func NewWorkerApp(cfg *config.AppConfigStruct) *WorkerApp {
	mqClient := queue.NewRabbitMQClient(&cfg.RabbitMQ)
	yaClient := crawler.NewYaMarketClient()
	consumer := worker.NewConsumer(mqClient, yaClient)

	return &WorkerApp{
		consumer: consumer,
		mq:       mqClient,
	}
}

// Run запускает потребление и завершает работу по SIGINT/SIGTERM.
func (a *WorkerApp) Run() error {
	defer func() { _ = a.mq.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() { errCh <- a.consumer.Run(ctx) }()

	select {
	case <-sig:
		log.Println("worker: получен сигнал, завершение...")
		cancel()
		return nil
	case err := <-errCh:
		return err
	}
}
