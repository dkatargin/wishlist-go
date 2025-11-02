package api

import (
	"wishlist-go/internal/api/handlers"
	"wishlist-go/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func PublicApi(router *gin.Engine) *gin.RouterGroup {
	publicEndpoints := router.Group("/api/v1/")
	publicEndpoints.Use(middleware.CorsMiddleware())
	publicEndpoints.Use(middleware.TelegramAuthMiddleware())
	{
		publicEndpoints.OPTIONS("*path", handlers.OptionsHandler)                      // Получить все списки
		publicEndpoints.GET("list", handlers.GetWishlists)                             // Получить все списки
		publicEndpoints.POST("list", handlers.CreateWishlist)                          // Создать список
		publicEndpoints.PATCH("list/:listId", handlers.UpdateWishlist)                 // Обновить список
		publicEndpoints.DELETE("list/:listId", handlers.DeleteWishlist)                // Удалить список
		publicEndpoints.GET("list/:listId/wishes", handlers.GetWishItems)              // Получить все желания в списке
		publicEndpoints.POST("list/:listId/wishes", handlers.CreateWishItem)           // Добавить желание в список
		publicEndpoints.GET("list/:listId/wishes/:wishId", handlers.GetWishItem)       // Получить конкретное желание
		publicEndpoints.PATCH("list/:listId/wishes/:wishId", handlers.UpdateWishItem)  // Обновить конкретное желание
		publicEndpoints.DELETE("list/:listId/wishes/:wishId", handlers.DeleteWishItem) // Удалить конкретное желание

		publicEndpoints.DELETE("account", handlers.DeleteAccount) // Удалить аккаунт и все списки

		publicEndpoints.GET("health", handlers.HealthCheck) // Проверка здоровья сервиса
	}
	return publicEndpoints
}
