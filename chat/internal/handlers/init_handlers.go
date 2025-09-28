package handlers

import (
	"microservices/chat/internal/svc/mongo_svc"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface {
	CreateRoomHandler(c *gin.Context)
}

type HandlerStruct struct {
	Mongo mongo_svc.MongoSvcInterface
}

func NewHandlers(
	mongo mongo_svc.MongoSvcInterface,
) *HandlerStruct {
	return &HandlerStruct{
		Mongo: mongo,
	}
}
