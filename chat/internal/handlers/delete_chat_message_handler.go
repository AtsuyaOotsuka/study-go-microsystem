package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"

	"github.com/gin-gonic/gin"
)

type DeleteChatMessageRequest struct {
	RoomID    string `form:"room_id" json:"room_id" binding:"required"`
	MessageID string `form:"message_id" json:"message_id" binding:"required"`
}

func (h *HandlerStruct) DeleteChatMessageHandler(c *gin.Context) {
	var req DeleteChatMessageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	roomID := req.RoomID
	userID := jwtinfo.UserID
	messageID := req.MessageID

	room, err := h.MongoSvc.GetRoomByID(roomID, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get room", "details": err.Error()})
		return
	}

	message, err := h.MongoSvc.GetChatMessageByID(roomID, messageID, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get message", "details": err.Error()})
		return
	}

	if message.UserID != int(userID) && room.OwnerID != int(userID) {
		c.JSON(403, gin.H{"error": "You can only delete your own messages"})
		return
	}

	err = h.MongoSvc.DeleteChatMessage(roomID, messageID, h.MongoPkg)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete message", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message deleted successfully"})
}
