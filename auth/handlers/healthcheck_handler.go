package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandlerInterface interface {
	Check(c *gin.Context)
}

type HealthCheckHandlerStruct struct{}

func NewHealthCheckHandler() *HealthCheckHandlerStruct {
	return &HealthCheckHandlerStruct{}
}

func (h *HealthCheckHandlerStruct) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
