package handler

import (
	"net/http"
	"wishlist-go/internal/delivery/http/dto"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.BasicResponse{
		Message: "ok",
	})
}

func OptionsHandler(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
