package model

import "time"

var ChatMessageCollectionName = "chat_messages"

type ChatMessage struct {
	RoomID        string
	UserID        int
	Message       string
	CreatedAt     time.Time
	IsReadUserIds []int
}
