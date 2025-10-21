package main

import (
	"wishlist-go/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	api.PublicApi(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
