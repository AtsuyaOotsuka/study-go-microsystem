package model

import "time"

var RoomCollectionName = "rooms"

type Room struct {
	Name      string
	OwnerID   int
	CreatedAt time.Time
}
