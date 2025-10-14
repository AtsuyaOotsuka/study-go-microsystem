package handlers

import (
	"microservices/chat/internal/svc/chat_svc"
	"microservices/chat/internal/svc/mongo_svc"
	"microservices/chat/pkg/mongo_pkg"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface {
	HealthCheckHandler(c *gin.Context)
	CreateRoomHandler(c *gin.Context)
	JoinRoomHandler(c *gin.Context)
	RoomListHandler(c *gin.Context)
	PostChatMessageHandler(c *gin.Context)
	LoadChatHandlers(c *gin.Context)
	ReadChatMessages(c *gin.Context)
	DeleteChatMessageHandler(c *gin.Context)
}

type HandlerStruct struct {
	MongoSvc mongo_svc.MongoSvcInterface
	MongoPkg mongo_pkg.MongoPkgInterface
	ChatSvc  chat_svc.ChatSvcInterface
}

func NewHandlers(
	mongo_svc mongo_svc.MongoSvcInterface,
	mongo_pkg mongo_pkg.MongoPkgInterface,
	chat_svc chat_svc.ChatSvcInterface,
) *HandlerStruct {
	return &HandlerStruct{
		MongoSvc: mongo_svc,
		MongoPkg: mongo_pkg,
		ChatSvc:  chat_svc,
	}
}
