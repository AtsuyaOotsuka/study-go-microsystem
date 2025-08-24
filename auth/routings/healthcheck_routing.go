package routings

import (
	"microservices/auth/handlers"

	"github.com/gin-gonic/gin"
)

func HealthCheckRouting(r *gin.Engine, handler handlers.HealthCheckHandlerInterface) {
	r.GET("/health", handler.Check)
}
