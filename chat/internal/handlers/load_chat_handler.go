package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *HandlerStruct) LoadChatHandlers(c *gin.Context) {
	roomID := c.Param("room_id")
	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	userID := jwtinfo.UserID
	room, err := h.MongoSvc.GetRoomByID(roomID, h.MongoPkg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room"})
		return
	}
	// ルームの情報を取得
	roomInfo := h.ChatSvc.GetRoomInfo(room, int(userID))
	if !roomInfo.IsMember && !roomInfo.IsOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	messages, err := h.MongoSvc.GetChatMessages(roomID, h.MongoPkg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": roomInfo, "messages": messages})
}
