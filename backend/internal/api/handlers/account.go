package handlers

import (
	"net/http"
	"wishlist-api/internal/api/middleware"
	"wishlist-api/internal/service"

	"github.com/gin-gonic/gin"
)

var accountService = service.NewAccountService()

func DeleteAccount(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID

	if err := accountService.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}
