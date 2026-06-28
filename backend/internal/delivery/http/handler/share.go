package handler

import (
	"errors"
	"net/http"
	"strconv"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/wishitem"
	"wishlist-go/internal/usecase/wishlist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShareHandler — публичный (без Telegram-auth) гостевой просмотр шаренного списка.
type ShareHandler struct {
	wishlistUC *wishlist.Service
	wishitemUC *wishitem.Service
}

func NewShareHandler(wlUC *wishlist.Service, wiUC *wishitem.Service) *ShareHandler {
	return &ShareHandler{wishlistUC: wlUC, wishitemUC: wiUC}
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

	resp := dto.SharedWishlistResponse{
		ShareCode:   wl.ShareCode,
		Name:        wl.Name,
		Description: wl.Description,
		Color:       wl.Color,
		Items:       make([]dto.SharedWishItem, 0, len(items)),
	}
	for _, it := range items {
		resp.Items = append(resp.Items, dto.SharedWishItem{
			ID:               it.ID,
			Name:             it.Name,
			Priority:         it.Priority,
			MarketURL:        it.MarketURL,
			MarketPictureURL: it.MarketPictureURL,
			MarketPrice:      it.MarketPrice,
			MarketCurrency:   it.MarketCurrency,
			MarketQuantity:   it.MarketQuantity,
		})
	}

	c.JSON(http.StatusOK, resp)
}
