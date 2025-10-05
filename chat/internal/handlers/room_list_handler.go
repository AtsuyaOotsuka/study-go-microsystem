package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"

	"github.com/gin-gonic/gin"
)

func (h *HandlerStruct) RoomListHandler(c *gin.Context) {
	target := c.DefaultQuery("target", "all")
	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	userID := jwtinfo.UserID

	rooms, err := h.MongoSvc.GetRooms(int(userID), target, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get rooms", "details": err.Error()})
		return
	}

	responseRooms := h.ChatSvc.ConvertRoomList(rooms, int(userID))

	c.JSON(200, gin.H{
		"rooms": responseRooms,
	})
}
