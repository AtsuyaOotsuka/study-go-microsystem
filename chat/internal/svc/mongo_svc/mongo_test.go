package mongo_svc

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"
	"microservices/chat/tests/mocks/pkg/mock_mongo_pkg"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInit(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{}, nil)

	_, err := Init(mongoPkgMock)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestInitError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	_, err := Init(mongoPkgMock)
	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestCreateRoom(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("mocked_id", nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
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
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	mockSvcStruct := NewMongoSvc(nil)

	room := model.Room{Name: "Test Room"}

	_, err := mockSvcStruct.CreateRoom(room, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestCreateRoomWithInsertError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionInsertErrorMock)
	mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("", assert.AnError)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
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

func TestGetRoomByID(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	var room model.Room
	mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &room).Return(nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	_, err := mockSvcStruct.GetRoomByID("64a7b2f4e13e4c3f9c8b4567", mongoPkgMock)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomByIDError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	mockSvcStruct := NewMongoSvc(nil)

	_, err := mockSvcStruct.GetRoomByID("64a7b2f4e13e4c3f9c8b4567", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomByIDWithObjectIDFromHexError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	_, err := mockSvcStruct.GetRoomByID("invalid_object_id", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomByIDWithFindOneError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	var room model.Room
	mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &room).Return(assert.AnError)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	_, err := mockSvcStruct.GetRoomByID("64a7b2f4e13e4c3f9c8b4567", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestJoinRoom(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoCollectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	err := mockSvcStruct.JoinRoom("64a7b2f4e13e4c3f9c8b4567", 1, mongoPkgMock)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestJoinRoomError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	mockSvcStruct := NewMongoSvc(nil)

	err := mockSvcStruct.JoinRoom("64a7b2f4e13e4c3f9c8b4567", 1, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestJoinRoomWithObjectIDFromHexError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	err := mockSvcStruct.JoinRoom("invalid_object_id", 1, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestJoinRoomWithUpdateOneError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoCollectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, assert.AnError)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)
	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	err := mockSvcStruct.JoinRoom("64a7b2f4e13e4c3f9c8b4567", 1, mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsForAll(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	var room model.Room
	mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)
	mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
	mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
	mongoCursorMock.On("Decode", &room).Return(nil)
	mongoCursorMock.On("Close", mock.Anything).Return(nil)

	filter := bson.M{
		"$or": []bson.M{
			{"isprivate": false},
			{"members": 1},
		},
	}

	mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
	_, err := mockSvcStruct.GetRooms(1, "all", mongoPkgMock)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsForJoined(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	var room model.Room
	mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)
	mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
	mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
	mongoCursorMock.On("Decode", &room).Return(nil)
	mongoCursorMock.On("Close", mock.Anything).Return(nil)

	filter := bson.M{"members": 1} // 参加済みのものだけ

	mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
	_, err := mockSvcStruct.GetRooms(1, "joined", mongoPkgMock)

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgWithErrorMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(nil, assert.AnError)

	mockSvcStruct := NewMongoSvc(nil)

	_, err := mockSvcStruct.GetRooms(1, "all", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsWithInvalidTarget(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

	_, err := mockSvcStruct.GetRooms(1, "invalid_target", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsWithFindError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)

	filter := bson.M{
		"$or": []bson.M{
			{"isprivate": false},
			{"members": 1},
		},
	}

	mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, assert.AnError)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
	_, err := mockSvcStruct.GetRooms(1, "all", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}

func TestGetRoomsWithDecodeError(t *testing.T) {
	mongoPkgMock := new(mock_mongo_pkg.MongoPkgMock)
	mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
	var room model.Room
	mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)
	mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
	mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
	mongoCursorMock.On("Decode", &room).Return(assert.AnError)
	mongoCursorMock.On("Close", mock.Anything).Return(nil)

	filter := bson.M{
		"$or": []bson.M{
			{"isprivate": false},
			{"members": 1},
		},
	}

	mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
	mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
	mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)
	mongoPkgMock.On("NewMongoConnect", "chatapp").Return(&mongo_pkg.MongoPkgStruct{
		Ctx:    context.TODO(),
		Db:     mongoDatabaseMock,
		Cancel: func() {},
	}, nil)

	mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
	_, err := mockSvcStruct.GetRooms(1, "all", mongoPkgMock)

	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	mongoPkgMock.AssertExpectations(t)
}
