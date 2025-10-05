package chat_svc

import "microservices/chat/internal/model"

type ChatSvcInterface interface {
	ConvertRoomList(rooms []model.Room, userId int) []Room
}

type ChatSvcStruct struct{}

func NewChatSvc() *ChatSvcStruct {
	return &ChatSvcStruct{}
}

type Room struct {
	ID          string
	Name        string
	OwnerID     int
	IsPrivate   bool
	IsMember    bool
	IsOwner     bool
	MemberCount int
	CreatedAt   string
}

func contains(members []int, target int) bool {
	for _, v := range members {
		if v == target {
			return true
		}
	}
	return false
}

func (s *ChatSvcStruct) ConvertRoomList(rooms []model.Room, userId int) []Room {
	var responses []Room
	for _, room := range rooms {
		response := Room{
			ID:          room.ID.Hex(),
			Name:        room.Name,
			OwnerID:     room.OwnerID,
			IsPrivate:   room.IsPrivate,
			IsMember:    contains(room.Members, userId),
			IsOwner:     room.OwnerID == userId,
			MemberCount: len(room.Members),
			CreatedAt:   room.CreatedAt.String(),
		}
		responses = append(responses, response)
	}
	return responses
}
