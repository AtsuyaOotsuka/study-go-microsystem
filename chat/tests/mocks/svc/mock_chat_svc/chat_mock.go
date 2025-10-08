package mock_chat_svc

import (
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/chat_svc"

	"github.com/stretchr/testify/mock"
)

type ChatSvcMock struct {
	mock.Mock
}

func (m *ChatSvcMock) ConvertRoomList(rooms []model.Room, userId int) []chat_svc.Room {
	args := m.Called(rooms, userId)
	return args.Get(0).([]chat_svc.Room)
}

func (m *ChatSvcMock) GetRoomInfo(room model.Room, userId int) chat_svc.Room {
	args := m.Called(room, userId)
	return args.Get(0).(chat_svc.Room)
}
