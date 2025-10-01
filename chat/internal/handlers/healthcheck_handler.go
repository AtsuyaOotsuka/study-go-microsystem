package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HandlerStruct) HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
