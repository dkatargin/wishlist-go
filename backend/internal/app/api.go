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
	"wishlist-go/internal/usecase/favorite"
	"wishlist-go/internal/usecase/reservation"
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
	reservationRepo := postgres.NewReservationRepository(db)
	favoriteRepo := postgres.NewFavoriteRepository(db)

	// Use cases
	accountUC := account.NewService(accountRepo)
	wishlistUC := wishlist.NewService(wishlistRepo)
	wishitemUC := wishitem.NewService(wishitemRepo, wishlistRepo, mqClient)
	reservationUC := reservation.NewService(reservationRepo, wishitemRepo)
	favoriteUC := favorite.NewService(favoriteRepo, wishlistRepo)

	// HTTP router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CorsMiddleware())
	// Handlers registration
	accountHandler := handler.NewAccountHandler(accountUC)
	wishlistHandler := handler.NewWishlistHandler(wishlistUC)
	wishitemHandler := handler.NewWishItemHandler(wishitemUC, wishlistUC)
	shareHandler := handler.NewShareHandler(wishlistUC, wishitemUC, reservationUC)
	reservationHandler := handler.NewReservationHandler(reservationUC)
	favoriteHandler := handler.NewFavoriteHandler(favoriteUC)
	// Routes
	api := router.Group("/api/v1")
	{
		api.OPTIONS("*path", handler.OptionsHandler)
		api.GET("health", handler.HealthCheck)
		// Гостевой просмотр шаренного списка: опциональная auth — пускаем и анонима,
		// но различаем владельца/гостя для правил видимости резервов.
		api.GET("share/:shareCode", middleware.OptionalTelegramAuth(cfg.Telegram.BotToken), shareHandler.Get)
		// Authorized routes
		authorized := api.Group("")
		authorized.Use(middleware.TelegramAuthMiddleware(cfg.Telegram.BotToken, accountUC))
		{
			authorized.GET("list", wishlistHandler.List)
			authorized.POST("list", wishlistHandler.Create)
			authorized.GET("list/:listId", wishlistHandler.Get)
			authorized.PATCH("list/:listId", wishlistHandler.Update)
			authorized.DELETE("list/:listId", wishlistHandler.Delete)
			authorized.GET("list/:listId/wishes", wishitemHandler.List)
			authorized.POST("list/:listId/wishes", wishitemHandler.Create)
			authorized.POST("list/:listId/wishes/crawl", wishitemHandler.Crawl)
			authorized.GET("list/:listId/wishes/:wishId", wishitemHandler.Get)
			authorized.PATCH("list/:listId/wishes/:wishId", wishitemHandler.Update)
			authorized.DELETE("list/:listId/wishes/:wishId", wishitemHandler.Delete)

			// Резервирование подарков (даритель — залогиненный гость)
			authorized.POST("share/:shareCode/wishes/:wishId/reserve", reservationHandler.Reserve)
			authorized.DELETE("share/:shareCode/wishes/:wishId/reserve", reservationHandler.Cancel)
			authorized.POST("reservations/:reservationId/purchased", reservationHandler.MarkPurchased)
			authorized.GET("reservations", reservationHandler.ListMine)

			// Избранное (сохранённые чужие списки); :id = ShareCode
			authorized.GET("favorites", favoriteHandler.List)
			authorized.POST("wishlist/:id/favorite", favoriteHandler.Add)
			authorized.DELETE("wishlist/:id/favorite", favoriteHandler.Remove)

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
