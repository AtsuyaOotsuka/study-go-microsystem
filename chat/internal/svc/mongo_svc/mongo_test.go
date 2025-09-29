package mongo_svc

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MongoPkgMock struct {
	mock.Mock
}

func (m *MongoPkgMock) NewMongoConnect(dbName string) (*mongo_pkg.MongoPkgStruct, error) {
	args := m.Called(dbName)
	return args.Get(0).(*mongo_pkg.MongoPkgStruct), args.Error(1)
}

type MongoPkgWithErrorMock struct {
	mock.Mock
}

func (m *MongoPkgWithErrorMock) NewMongoConnect(dbName string) (*mongo_pkg.MongoPkgStruct, error) {
	args := m.Called(dbName)
	return nil, args.Error(1)
}

func TestInit(t *testing.T) {
	mongoPkgMock := new(MongoPkgMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{}, nil)

	_, err := Init(mongoPkgMock)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestInitError(t *testing.T) {
	mongoPkgMock := new(MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	_, err := Init(mongoPkgMock)
	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

type MongoDatabaseMock struct {
	mock.Mock
}

func (m *MongoDatabaseMock) Collection(name string) mongo_pkg.MongoCollectionInterface {
	args := m.Called(name)
	return args.Get(0).(mongo_pkg.MongoCollectionInterface)
}

type MongoCollectionMock struct {
	mock.Mock
}

func (m *MongoCollectionMock) InsertOne(ctx context.Context, document interface{}) (string, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(string), args.Error(1)
}

type MongoPkgStructMock struct {
	Ctx context.Context
	Db  mongo_pkg.MongoDatabaseInterface
}

type MongoMock struct {
	MongoPkgStruct *MongoPkgStructMock
}

func TestCreateRoom(t *testing.T) {
	mongoPkgMock := new(MongoPkgMock)
	mongoCollectionMock := new(MongoCollectionMock)
	mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("mocked_id", nil)
	mongoDatabaseMock := new(MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	room := model.Room{Name: "Test Room"}

	_, err := mockSvcStruct.CreateRoom(room, mongoPkgMock)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestCreateRoomError(t *testing.T) {
	mongoPkgMock := new(MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	mockSvcStruct := NewMongoSvc(nil)

	room := model.Room{Name: "Test Room"}

	_, err := mockSvcStruct.CreateRoom(room, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

type MongoCollectionInsertErrorMock struct {
	mock.Mock
}

func (m *MongoCollectionInsertErrorMock) InsertOne(ctx context.Context, document interface{}) (string, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(string), args.Error(1)
}

func TestCreateRoomWithInsertError(t *testing.T) {
	mongoPkgMock := new(MongoPkgMock)
	mongoCollectionMock := new(MongoCollectionInsertErrorMock)
	mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("", assert.AnError)
	mongoDatabaseMock := new(MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	room := model.Room{Name: "Test Room"}

	_, err := mockSvcStruct.CreateRoom(room, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}
