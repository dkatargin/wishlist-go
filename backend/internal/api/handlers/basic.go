package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func OptionsHandler(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
