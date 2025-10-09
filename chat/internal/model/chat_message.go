package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ChatMessageCollectionName = "chat_messages"

type ChatMessage struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	RoomID        string
	UserID        int
	Message       string
	CreatedAt     time.Time
	IsReadUserIds []int
}
