package handlers

import (
	"microservices/chat/internal/svc/mongo_svc"
	"microservices/chat/pkg/mongo_pkg"

	"github.com/gin-gonic/gin"
)

type HandlersInterface interface {
	HealthCheckHandler(c *gin.Context)
	CreateRoomHandler(c *gin.Context)
}

type HandlerStruct struct {
	MongoSvc mongo_svc.MongoSvcInterface
	MongoPkg mongo_pkg.MongoPkgInterface
}

func NewHandlers(
	mongo_svc mongo_svc.MongoSvcInterface,
	mongo_pkg mongo_pkg.MongoPkgInterface,
) *HandlerStruct {
	return &HandlerStruct{
		MongoSvc: mongo_svc,
		MongoPkg: mongo_pkg,
	}
}
