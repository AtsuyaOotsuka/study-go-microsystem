package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"

	"github.com/gin-gonic/gin"
)

type JoinRoomRequest struct {
	RoomID string `form:"room_id" json:"room_id" binding:"required"`
}

func (h *HandlerStruct) JoinRoomHandler(c *gin.Context) {
	var req JoinRoomRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	roomID := req.RoomID
	userID := jwtinfo.UserID

	_, err := h.MongoSvc.GetRoomByID(roomID, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get room", "details": err.Error()})
		return
	}

	err = h.MongoSvc.JoinRoom(roomID, int(userID), h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to join room", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Joined room successfully"})
}
