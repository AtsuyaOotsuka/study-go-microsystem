package routings

import (
	"microservices/chat/internal/handlers"

	"github.com/gin-gonic/gin"
)

type Rootings struct {
	CSRFMW gin.HandlerFunc
	AuthMW gin.HandlerFunc
}

func Routing(r *gin.Engine, csrfMW gin.HandlerFunc, authMW gin.HandlerFunc, handlers handlers.HandlersInterface) {
	r.Use(csrfMW, authMW)
	r.POST("/room_create", handlers.CreateRoomHandler)
	r.GET("/health", handlers.HealthCheckHandler)
}
