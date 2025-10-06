package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"

	"github.com/gin-gonic/gin"
)

type PostChatRequest struct {
	RoomID  string `form:"room_id" json:"room_id" binding:"required"`
	Message string `form:"message" json:"message" binding:"required"`
}

func (h *HandlerStruct) PostChatMessageHandler(c *gin.Context) {
	var req PostChatRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	roomID := req.RoomID
	userID := jwtinfo.UserID
	message := req.Message

	_, err := h.MongoSvc.GetRoomByID(roomID, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get room", "details": err.Error()})
		return
	}

	err = h.MongoSvc.PostChatMessage(roomID, int(userID), message, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to post chat", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Chat posted successfully"})
}
