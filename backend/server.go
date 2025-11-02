package main

import (
	"flag"
	"log"
	"wishlist-go/internal/api"
	"wishlist-go/internal/config"
	"wishlist-go/internal/db"
	"wishlist-go/internal/queue"

	"github.com/gin-gonic/gin"
)

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
		panic("Failed to connect to the database: " + err.Error())
	}

	// Подключаемся к RabbitMQ
	err = queue.ConnectRabbitMQ()
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		log.Println("Continuing without RabbitMQ support")
	} else {
		defer queue.Client.Close()
	}

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api.PublicApi(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
