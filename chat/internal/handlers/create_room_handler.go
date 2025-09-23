package handlers

import (
	"fmt"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateRoomRequest struct {
	Name string `form:"name" json:"name" binding:"required"`
}

func (h *HandlerStruct) CreateRoomHandler(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	h.Mongo.ReConnect()
	defer h.Mongo.Cancel()
	collection := h.Mongo.Db.Collection(model.RoomCollectionName)

	room := model.Room{
		Name:      req.Name,
		OwnerID:   jwtinfo.UserID,
		CreatedAt: time.Now(),
	}

	result, err := collection.InsertOne(h.Mongo.Ctx, room)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create room", "details": err.Error()})
		return
	}
	fmt.Print(result)

	c.JSON(200, gin.H{
		"message":    "Room created successfully",
		"room_id":    result.InsertedID,
		"room_name":  room.Name,
		"created_at": room.CreatedAt,
	})
}
