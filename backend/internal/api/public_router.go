package api

import (
	"wishlist-go/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func PublicApi(router *gin.Engine) *gin.RouterGroup {
	publicEndpoints := router.Group("/api/v1/")
	publicEndpoints.Use(handlers.TelegramAuthMiddleware())
	{
		publicEndpoints.GET("list", handlers.BasicHandler)                           // Получить все списки
		publicEndpoints.POST("list", handlers.BasicHandler)                          // Создать список
		publicEndpoints.PATCH("list/:listId", handlers.BasicHandler)                 // Обновить список
		publicEndpoints.DELETE("list/:listId", handlers.BasicHandler)                // Удалить список
		publicEndpoints.GET("list/:listId/wishes", handlers.BasicHandler)            // Получить все желания в списке
		publicEndpoints.POST("list/:listId/wishes", handlers.BasicHandler)           // Добавить желание в список
		publicEndpoints.POST("list/:listId/wishes", handlers.BasicHandler)           // Создать желание
		publicEndpoints.GET("list/:listId/wishes/:wishId", handlers.BasicHandler)    // Получить конкретное желание
		publicEndpoints.PATCH("list/:listId/wishes/:wishId", handlers.BasicHandler)  // Обновить конкретное желание
		publicEndpoints.DELETE("list/:listId/wishes/:wishId", handlers.BasicHandler) // Удалить конкретное желание

		publicEndpoints.DELETE("account", handlers.BasicHandler) // Удалить аккаунт и все списки
	}
	return publicEndpoints
}
