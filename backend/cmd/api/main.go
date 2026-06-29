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

	apiApp := app.NewAPIApp(cfg)
	if err := apiApp.Run(); err != nil {
		log.Fatal(err)
	}
}
