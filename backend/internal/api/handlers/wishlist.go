package handlers

import (
	"net/http"
	"strconv"
	"wishlist-go/internal/api/middleware"
	"wishlist-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetWishlists(c *gin.Context) {
	var wishlistService = service.NewWishlistService()
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid offset"})
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid limit"})
		return
	}

	owner := &auth.(*middleware.TelegramAuthData).User.ID

	wishlists, err := wishlistService.GetAllByOwner(*owner, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error fetching wishlists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wishlists": wishlists})
}

func CreateWishlist(c *gin.Context) {
	var wishlistService = service.NewWishlistService()
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	type req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	wl := service.WishlistInsert{
		Owner:       &auth.(*middleware.TelegramAuthData).User.ID,
		Name:        &r.Name,
		Description: &r.Description,
	}

	wishlist, err := wishlistService.Create(&wl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}

func UpdateWishlist(c *gin.Context) {
	var wishlistService = service.NewWishlistService()
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}

	// Проверяем, есть ли у пользователя доступ к этому списку желаний
	hasAccess := wishlistService.CheckAccess(auth.(*middleware.TelegramAuthData).User.ID, shareCode)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	type req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	wishlist, err := wishlistService.Update(shareCode, service.WishlistInsert{
		Name:        r.Name,
		Description: r.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wishlist": wishlist})
}

func DeleteWishlist(c *gin.Context) {
	var wishlistService = service.NewWishlistService()
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}
	// Проверяем, есть ли у пользователя доступ к этому списку желаний
	hasAccess := wishlistService.CheckAccess(auth.(*middleware.TelegramAuthData).User.ID, shareCode)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	err = wishlistService.Delete(shareCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting wishlist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "wishlist deleted"})
}
