package routings

import (
	"microservices/auth/handlers"

	"github.com/gin-gonic/gin"
)

func CsrfRouting(r *gin.Engine, handler handlers.CSRFHandlerInterface) {
	routerGroup := r.Group("/csrf")
	routerGroup.GET("/get", handler.CsrfGet)
}
