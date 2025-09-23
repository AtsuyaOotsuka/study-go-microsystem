package handlers

import (
	"microservices/chat/pkg/mongo_pkg"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface {
	CreateRoomHandler(c *gin.Context)
}

type HandlerStruct struct {
	Mongo *mongo_pkg.MongoPkgStruct
}

func NewHandlers(mongo *mongo_pkg.MongoPkgStruct) *HandlerStruct {
	return &HandlerStruct{
		Mongo: mongo,
	}
}
