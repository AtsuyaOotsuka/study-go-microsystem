package routings

import (
	"microservices/auth/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRouting(r *gin.Engine, handlerFunc handlers.RegisterHandlerInterface, csrfMW gin.HandlerFunc) {
	r.Use(csrfMW)
	r.POST("/register", handlerFunc.HandleRegister)
}
