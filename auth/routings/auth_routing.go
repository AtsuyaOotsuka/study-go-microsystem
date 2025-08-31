package routings

import (
	"microservices/auth/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRouting(r *gin.Engine, handlerFunc handlers.AuthHandlerInterface, csrfMW gin.HandlerFunc) {
	routerGroup := r.Group("/auth")
	routerGroup.Use(csrfMW)
	routerGroup.POST("/login", handlerFunc.HandleLogin)
	routerGroup.POST("/refresh", handlerFunc.HandleRefresh)
}
