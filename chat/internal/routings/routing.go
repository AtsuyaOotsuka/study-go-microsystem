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
	r.POST("/room_join", handlers.JoinRoomHandler)
	r.GET("/room_list", handlers.RoomListHandler)
	r.POST("/post_chat_message", handlers.PostChatMessageHandler)
	r.GET("/health", handlers.HealthCheckHandler)
}
