package handler

import (
	"errors"
	"net/http"
	"strconv"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/reservation"
	"wishlist-go/internal/usecase/wishitem"
	"wishlist-go/internal/usecase/wishlist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShareHandler — гостевой просмотр шаренного списка под OptionalTelegramAuth:
// пускает и анонима, но различает владельца / гостя / анонима для правил видимости резервов.
type ShareHandler struct {
	wishlistUC    *wishlist.Service
	wishitemUC    *wishitem.Service
	reservationUC *reservation.Service
}

func NewShareHandler(wlUC *wishlist.Service, wiUC *wishitem.Service, resUC *reservation.Service) *ShareHandler {
	return &ShareHandler{wishlistUC: wlUC, wishitemUC: wiUC, reservationUC: resUC}
}

// remainingQuantity — сколько ещё можно зарезервировать.
// nil, если market_quantity не задан (0 трактуем как «не указано» — лимита нет);
// перебронь не уходит в минус.
func remainingQuantity(marketQuantity *int, reserved int) *int {
	if marketQuantity == nil || *marketQuantity <= 0 {
		return nil
	}
	rem := *marketQuantity - reserved
	if rem < 0 {
		rem = 0
	}
	return &rem
}

// applyReservationView заполняет поля резерва согласно правам смотрящего:
//   - владелец списка не видит резервы вообще (сюрприз сохраняется);
//   - аноним/неаутентифицированный видит только агрегат количества, без состава дарителей;
//   - аутентифицированный гость дополнительно видит не-анонимных дарителей.
//
// Reservers всегда не-nil (пустой срез, а не null).
func applyReservationView(si *dto.SharedWishItem, isOwner, isAuthenticated bool, agg reservation.WishReservationView) {
	si.Reservers = []int64{}
	if isOwner {
		return
	}
	si.ReservedQuantity = agg.ReservedQuantity
	si.Remaining = remainingQuantity(si.MarketQuantity, agg.ReservedQuantity)
	if isAuthenticated && len(agg.Reservers) > 0 {
		si.Reservers = agg.Reservers
	}
}

// Get отдаёт список (публичные поля) и его невыполненные желания по share-коду.
func (h *ShareHandler) Get(c *gin.Context) {
	shareCode, err := uuid.Parse(c.Param("shareCode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid share code"})
		return
	}

	wl, err := h.wishlistUC.GetByCode(c.Request.Context(), shareCode)
	if err != nil {
		if errors.Is(err, domain.ErrWishlistNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wishlist not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wishlist"})
		return
	}

	limit := 50
	if v, err := strconv.Atoi(c.DefaultQuery("limit", "50")); err == nil && v > 0 {
		if v > 100 {
			v = 100
		}
		limit = v
	}
	offset := 0
	if v, err := strconv.Atoi(c.DefaultQuery("offset", "0")); err == nil && v >= 0 {
		offset = v
	}
	items, err := h.wishitemUC.GetVisibleByWishlist(c.Request.Context(), shareCode, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wish items"})
		return
	}

	// Кто смотрит (OptionalTelegramAuth кладёт telegram_auth только при валидной подписи).
	var requesterID int64
	if auth, ok := c.Get("telegram_auth"); ok {
		requesterID = auth.(*middleware.TelegramAuthData).User.ID
	}
	isOwner := requesterID == wl.OwnerID
	isAuthenticated := requesterID != 0

	// Владельцу резервы не показываем вообще — и не дёргаем агрегат лишний раз.
	var reservations map[int64]reservation.WishReservationView
	if !isOwner {
		reservations, err = h.reservationUC.AggregateForSharedList(c.Request.Context(), shareCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch reservations"})
			return
		}
	}

	resp := dto.SharedWishlistResponse{
		ShareCode:   wl.ShareCode,
		Name:        wl.Name,
		Description: wl.Description,
		Color:       wl.Color,
		Items:       make([]dto.SharedWishItem, 0, len(items)),
	}
	for _, it := range items {
		si := dto.SharedWishItem{
			ID:               it.ID,
			Name:             it.Name,
			Priority:         it.Priority,
			MarketURL:        it.MarketURL,
			MarketPictureURL: it.MarketPictureURL,
			MarketPrice:      it.MarketPrice,
			MarketCurrency:   it.MarketCurrency,
			MarketQuantity:   it.MarketQuantity,
		}
		applyReservationView(&si, isOwner, isAuthenticated, reservations[it.ID])
		resp.Items = append(resp.Items, si)
	}

	c.JSON(http.StatusOK, resp)
}
