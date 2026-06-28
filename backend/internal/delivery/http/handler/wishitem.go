package handler

import (
	"errors"
	"net/http"
	"strconv"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/domain"
	"wishlist-go/internal/usecase/wishitem"
	"wishlist-go/internal/usecase/wishlist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishItemHandler struct {
	usecase         *wishitem.Service
	wishlistUsecase *wishlist.Service
}

func NewWishItemHandler(uc *wishitem.Service, wlUC *wishlist.Service) *WishItemHandler {
	return &WishItemHandler{
		usecase:         uc,
		wishlistUsecase: wlUC,
	}
}

func (h *WishItemHandler) Create(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}

	var req dto.CreateWishItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	// Проверяем доступ пользователя к списку желаний
	hasAccess := h.wishlistUsecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	item := &domain.WishItem{
		WishListCode:     shareCode,
		OwnerID:          userID,
		Name:             &req.Name,
		Priority:         req.Priority,
		MarketURL:        req.MarketLink,
		MarketPictureURL: &req.MarketPicture,
		MarketPrice:      &req.MarketPrice,
		MarketCurrency:   req.MarketCurrency,
		MarketQuantity:   &req.MarketQuantity,
	}
	wi, err := h.usecase.CreateWishItem(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating wish item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wish_item": wi})

}

func (h *WishItemHandler) Get(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}

	itemId := c.Param("wishId")
	wishId, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid item id"})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	// Проверяем доступ пользователя к списку желаний
	hasAccess := h.wishlistUsecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	wi, err := h.usecase.GetWishItemByID(c.Request.Context(), wishId, shareCode)
	if err != nil {
		if errors.Is(err, domain.ErrWishItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wish item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching wish item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wish_item": wi})
}

func (h *WishItemHandler) List(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "50")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	// Проверяем доступ пользователя к списку желаний
	hasAccess := h.wishlistUsecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	wishitems, err := h.usecase.GetWishItemsByWishlist(c.Request.Context(), shareCode, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wish items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wish_items": wishitems})

}

func (h *WishItemHandler) Update(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}

	itemId := c.Param("wishId")
	wishId, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid item id"})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	// Проверяем доступ пользователя к списку желаний
	hasAccess := h.wishlistUsecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req dto.UpdateWishItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	upd := domain.WishItemUpdate{
		Name:             req.Name,
		Priority:         req.Priority,
		IsDone:           req.IsDone,
		MarketURL:        req.MarketLink,
		MarketPictureURL: req.MarketPicture,
		MarketPrice:      req.MarketPrice,
		MarketCurrency:   req.MarketCurrency,
		MarketQuantity:   req.MarketQuantity,
	}
	if err := h.usecase.UpdateWishItem(c.Request.Context(), wishId, shareCode, upd); err != nil {
		if errors.Is(err, domain.ErrWishItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wish item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating wish item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}

func (h *WishItemHandler) Delete(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}

	itemId := c.Param("wishId")
	wishId, err := strconv.ParseInt(itemId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid item id"})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	// Проверяем доступ пользователя к списку желаний
	hasAccess := h.wishlistUsecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = h.usecase.DeleteWishItem(c.Request.Context(), wishId)
	if err != nil {
		if errors.Is(err, domain.ErrWishItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wish item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting wish item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
