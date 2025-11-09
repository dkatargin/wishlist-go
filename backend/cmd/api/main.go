package main

import (
	"flag"
	"log"
	"wishlist-go/internal/app"
	"wishlist-go/internal/infrastructure/config"
)

func main() {
	cfgPath := flag.String("config", "configs/config.yaml", "config path")
	flag.Parse()

	cfg, err := config.LoadConfigFile(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	apiApp := app.NewAPIApp(cfg)
	if err := apiApp.Run(); err != nil {
		log.Fatal(err)
	}

}
