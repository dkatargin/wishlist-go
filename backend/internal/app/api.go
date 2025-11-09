package app

import (
	"fmt"
	"wishlist-go/internal/delivery/http/handler"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/infrastructure/config"
	"wishlist-go/internal/infrastructure/database"
	"wishlist-go/internal/infrastructure/queue"
	"wishlist-go/internal/repository/postgres"
	"wishlist-go/internal/usecase/account"
	"wishlist-go/internal/usecase/wishitem"
	"wishlist-go/internal/usecase/wishlist"

	"github.com/gin-gonic/gin"
)

type APIApp struct {
	router *gin.Engine
	cfg    *config.AppConfigStruct
}

func NewAPIApp(cfg *config.AppConfigStruct) *APIApp {
	// Infrastructure initialization
	db := database.ConnectDB(&cfg.Database)
	mqClient := queue.NewRabbitMQClient(&cfg.RabbitMQ)

	// Repository
	accountRepo := postgres.NewAccountRepository(db)
	wishlistRepo := postgres.NewWishlistRepository(db)
	wishitemRepo := postgres.NewWishItemRepository(db)

	// Use cases
	accountUC := account.NewService(accountRepo)
	wishlistUC := wishlist.NewService(wishlistRepo)
	wishitemUC := wishitem.NewService(wishitemRepo, wishlistRepo, mqClient)

	// HTTP router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CorsMiddleware())
	// Handlers registration
	accountHandler := handler.NewAccountHandler(accountUC)
	wishlistHandler := handler.NewWishlistHandler(wishlistUC)
	wishitemHandler := handler.NewWishItemHandler(wishitemUC)
	// Routes
	api := router.Group("/api/v1")
	{
		api.OPTIONS("*path", handler.OptionsHandler)
		api.GET("health", handler.HealthCheck)
		// Authorized routes
		authorized := api.Group("")
		authorized.Use(middleware.TelegramAuthMiddleware())
		{
			authorized.GET("list", wishlistHandler.List)
			authorized.POST("list", wishlistHandler.Create)
			authorized.PATCH("list/:listId", wishlistHandler.Update)
			authorized.DELETE("list/:listId", wishlistHandler.Delete)
			authorized.GET("list/:listId/wishes", wishitemHandler.List)
			authorized.POST("list/:listId/wishes", wishitemHandler.Create)
			authorized.GET("list/:listId/wishes/:wishId", wishitemHandler.Get)
			authorized.PATCH("list/:listId/wishes/:wishId", wishitemHandler.Update)
			authorized.DELETE("list/:listId/wishes/:wishId", wishitemHandler.Delete)

			authorized.DELETE("account", accountHandler.Delete)
		}
	}
	return &APIApp{
		router: router,
		cfg:    cfg,
	}
}

func (a *APIApp) Run() error {
	return a.router.Run(fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port))
}
