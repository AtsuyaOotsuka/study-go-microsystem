package mock_mongo_svc

import (
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"

	"github.com/stretchr/testify/mock"
)

type MongoSvcMock struct {
	mock.Mock
}

func (m *MongoSvcMock) CreateRoom(room model.Room, mongo_pkg mongo_pkg.MongoPkgInterface) (interface{}, error) {
	args := m.Called(room, mongo_pkg)
	return args.Get(0), args.Error(1)
}

func (m *MongoSvcMock) GetRoomByID(roomID string, mongo_pkg mongo_pkg.MongoPkgInterface) (model.Room, error) {
	args := m.Called(roomID, mongo_pkg)
	return args.Get(0).(model.Room), args.Error(1)
}

func (m *MongoSvcMock) JoinRoom(roomID string, userID int, mongo_pkg mongo_pkg.MongoPkgInterface) error {
	args := m.Called(roomID, userID, mongo_pkg)
	return args.Error(0)
}

type MongoSvcMockWithErrorMock struct {
	mock.Mock
}

func (m *MongoSvcMockWithErrorMock) CreateRoom(room model.Room, mongo_pkg mongo_pkg.MongoPkgInterface) (interface{}, error) {
	args := m.Called(room, mongo_pkg)
	return args.Get(0), args.Error(1)
}
func (m *MongoSvcMockWithErrorMock) GetRoomByID(roomID string, mongo_pkg mongo_pkg.MongoPkgInterface) (model.Room, error) {
	args := m.Called(roomID, mongo_pkg)
	return args.Get(0).(model.Room), args.Error(1)
}
func (m *MongoSvcMockWithErrorMock) JoinRoom(roomID string, userID int, mongo_pkg mongo_pkg.MongoPkgInterface) error {
	args := m.Called(roomID, userID, mongo_pkg)
	return args.Error(0)
}
