package handler

import (
	"net/http"
	"wishlist-go/internal/delivery/http/dto"
	"wishlist-go/internal/delivery/http/middleware"
	"wishlist-go/internal/usecase/account"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	usecase *account.AccountService
}

func NewAccountHandler(e *account.AccountService) *AccountHandler {
	return &AccountHandler{
		usecase: e,
	}
}

func (h *AccountHandler) Delete(c *gin.Context) {
	auth, exist := c.Get("telegram_auth")
	if !exist {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	userID := auth.(*middleware.TelegramAuthData).User.ID

	if err := h.usecase.Delete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.BasicResponse{Message: "account deleted"})
}
