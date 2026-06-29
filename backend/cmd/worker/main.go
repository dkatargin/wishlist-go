package main

import (
	"log"
	"wishlist-go/internal/app"
	"wishlist-go/internal/infrastructure/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	workerApp := app.NewWorkerApp(cfg)
	if err := workerApp.Run(); err != nil {
		log.Fatal(err)
	}
}
