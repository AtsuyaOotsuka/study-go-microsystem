package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var RoomCollectionName = "rooms"

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string
	OwnerID   int
	CreatedAt time.Time
	Members   []int
	IsPrivate bool
}
