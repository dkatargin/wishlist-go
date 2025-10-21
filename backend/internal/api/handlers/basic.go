package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "Hello World"})
}
