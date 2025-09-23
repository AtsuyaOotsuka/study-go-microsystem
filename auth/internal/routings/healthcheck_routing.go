package routings

import (
	"microservices/auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func HealthCheckRouting(r *gin.Engine, handler handlers.HealthCheckHandlerInterface) {
	r.GET("/health", handler.Check)
}
