package handlers

import (
	"microservices/chat/internal/svc/jwtinfo_svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReadChatRequest struct {
	RoomID     string   `form:"room_id" json:"room_id" binding:"required"`
	ChatIDList []string `form:"chat_id_list" json:"chat_id_list" binding:"required"`
}

func (h *HandlerStruct) ReadChatMessages(c *gin.Context) {
	var req ReadChatRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	jwtinfo := jwtinfo_svc.NewJwtInfo(c.Request.Context())
	userID := jwtinfo.UserID
	room, err := h.MongoSvc.GetRoomByID(req.RoomID, h.MongoPkg)
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
	chatIdList := req.ChatIDList
	if len(chatIdList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id_list is required"})
		return
	}
	err = h.MongoSvc.ReadChatMessages(req.RoomID, chatIdList, int(userID), h.MongoPkg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read chat message", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat messages marked as read"})

}
