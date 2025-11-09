package handler

import (
	"net/http"
	"strconv"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/usecase/wishlist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	usecase *wishlist.Service
}

func NewWishlistHandler(uc *wishlist.Service) *WishlistHandler {
	return &WishlistHandler{
		usecase: uc,
	}

}

func (h *WishlistHandler) Create(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := auth.(*middleware.TelegramAuthData).User.ID

	wl, err := h.usecase.CreateWishlist(c.Request.Context(), userID, req.Name, &req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create wishlist"})
		return
	}

	c.JSON(http.StatusOK, wl)
}

func (h *WishlistHandler) List(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "10")

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

	wishlists, err := h.usecase.GetWishlistsByOwner(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch wishlists"})
		return
	}

	c.JSON(http.StatusOK, wishlists)
}

func (h *WishlistHandler) Update(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	userID := auth.(*middleware.TelegramAuthData).User.ID
	var req dto.UpdateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}

	// Проверяем, есть ли у пользователя доступ к этому списку желаний
	hasAccess := h.usecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = h.usecase.UpdateWishlist(c.Request.Context(), shareCode, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{"error": "unauthorized"})
		return
	}

	userID := auth.(*middleware.TelegramAuthData).User.ID
	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}

	// Проверяем, есть ли у пользователя доступ к этому списку желаний
	hasAccess := h.usecase.CheckAccess(c.Request.Context(), shareCode, userID)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = h.usecase.DeleteWishlist(c.Request.Context(), shareCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})

}
