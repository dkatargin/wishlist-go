package handlers

import (
	"net/http"
	"strconv"
	"wishlist-go/internal/api/middleware"
	"wishlist-go/internal/queue"
	"wishlist-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var wishItemService = service.NewWishItemService()

func GetWishItems(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	listId := c.Param("listId")
	shareCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid list id"})
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid limit"})
		return
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid offset"})
		return
	}

	wishItemService.Owner = &auth.(*middleware.TelegramAuthData).User.ID
	wishItemService.WishList = shareCode
	wishItems, err := wishItemService.GetAll(limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "error fetching wish items"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wish_items": wishItems})

}

func CreateWishItem(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	listId := c.Param("listId")
	wishListCode, err := uuid.Parse(listId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid list id"})
		return
	}

	wlService := service.NewWishlistService()
	hasAccess := wlService.CheckAccess(auth.(*middleware.TelegramAuthData).User.ID, wishListCode)
	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	type req struct {
		Name           string  `json:"name" binding:"required"`
		Priority       int     `json:"priority"`
		MarketLink     string  `json:"market_link"`
		MarketPicture  string  `json:"market_picture"`
		MarketPrice    float64 `json:"market_price"`
		MarketCurrency string  `json:"market_currency"`
		MarketQuantity int     `json:"market_quantity"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	defaultStatus := "pending"

	wishItemService.WishList = wishListCode
	wishItemService.Owner = &auth.(*middleware.TelegramAuthData).User.ID
	wishItemService.Name = &r.Name
	wishItemService.Priority = &r.Priority
	wishItemService.Status = &defaultStatus // pending by default
	wishItemService.MarketLink = &r.MarketLink
	wishItemService.MarketPicture = &r.MarketPicture
	wishItemService.MarketPrice = &r.MarketPrice
	wishItemService.MarketCurrency = &r.MarketCurrency
	wishItemService.MarketQuantity = &r.MarketQuantity
	wishItem, err := wishItemService.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating wish item"})
		return
	}

	// Отправляем сообщение в RabbitMQ о создании wishitem
	if queue.Client != nil {
		_ = queue.Client.PublishMessage("wishitem_created", map[string]interface{}{
			"wishitem_id":  wishItem.ID,
			"wishlist_id":  wishListCode.String(),
			"owner_id":     *wishItemService.Owner,
			"name":         *wishItemService.Name,
			"market_price": *wishItemService.MarketPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{"wish_item": wishItem})
}

func GetWishItem(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wishItemID := c.Param("wishId")
	id, err := strconv.ParseInt(wishItemID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish item id"})
		return
	}

	wishItemService.Owner = &auth.(*middleware.TelegramAuthData).User.ID
	wishItem, err := wishItemService.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching wish item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wish_item": wishItem})
}

func UpdateWishItem(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wishItemID := c.Param("wishId")
	id, err := strconv.ParseInt(wishItemID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish item id"})
		return
	}

	type req struct {
		Name           *string  `json:"name"`
		Priority       *int     `json:"priority"`
		Status         *string  `json:"status"`
		MarketLink     *string  `json:"market_link"`
		MarketPicture  *string  `json:"market_picture"`
		MarketPrice    *float64 `json:"market_price"`
		MarketCurrency *string  `json:"market_currency"`
		MarketQuantity *int     `json:"market_quantity"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	wishItemService.Owner = &auth.(*middleware.TelegramAuthData).User.ID
	wishItemService.Name = r.Name
	wishItemService.Priority = r.Priority
	wishItemService.Status = r.Status
	wishItemService.MarketLink = r.MarketLink
	wishItemService.MarketPicture = r.MarketPicture
	wishItemService.MarketPrice = r.MarketPrice
	wishItemService.MarketCurrency = r.MarketCurrency
	wishItemService.MarketQuantity = r.MarketQuantity

	wishItem, err := wishItemService.Update(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating wish item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wish_item": wishItem})
}

func DeleteWishItem(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	wishItemID := c.Param("wishId")
	id, err := strconv.ParseInt(wishItemID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wish item id"})
		return
	}

	wishItemService.Owner = &auth.(*middleware.TelegramAuthData).User.ID
	err = wishItemService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting wish item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "wish item deleted"})
}
