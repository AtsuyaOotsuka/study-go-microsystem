package routings

import (
	"microservices/auth/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRouting(r *gin.Engine, handlerFunc handlers.RegisterHandlerInterface, csrfMW gin.HandlerFunc) {
	r.Use(csrfMW)
	r.POST("/register", handlerFunc.Register)
}
