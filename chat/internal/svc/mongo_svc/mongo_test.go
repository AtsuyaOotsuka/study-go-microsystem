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

func setupInitMock(wantErr bool, collection string, returnVal interface{}) mongo_pkg.MongoPkgInterface {
	if wantErr {
		m := new(mock_mongo_pkg.MongoPkgWithErrorMock)
		m.On("NewMongoConnect", collection).Return(nil, assert.AnError)
		return m
	}
	m := new(mock_mongo_pkg.MongoPkgMock)
	m.On("NewMongoConnect", collection).Return(returnVal, nil)
	return m
}

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		initErr bool
	}{
		{"success", false},
		{"error", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := setupInitMock(tt.initErr, "chatapp", &mongo_pkg.MongoPkgStruct{})
			_, err := Init(mock)
			if m, ok := mock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
			if (err != nil) != tt.initErr {
				t.Errorf("Init() [%s] error = %v, wantErr %v", tt.name, err, tt.initErr)
			}
		})
	}
}

func TestCreateRoom(t *testing.T) {
	tests := []struct {
		name         string
		initErr      bool
		InsertOneErr bool
	}{
		{"success", false, false},
		{"error", true, false},
		{"insert_error", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			if tt.InsertOneErr {
				mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("", assert.AnError)
			} else {
				mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("mocked_id", nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)
			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
			room := model.Room{Name: "Test Room"}

			_, err := mockSvcStruct.CreateRoom(room, mongoPkgMock)
			if (err != nil) != tt.initErr && (err != nil) != tt.InsertOneErr {
				t.Errorf("CreateRoom() [%s] error = %v, wantErr %v", tt.name, err, tt.initErr)
			}
			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestGetRoomByID(t *testing.T) {
	tests := []struct {
		name       string
		initErr    bool
		request    string
		findOneErr bool
		returnErr  bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4567", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4567", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"findone_error", false, "64a7b2f4e13e4c3f9c8b4567", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			var room model.Room
			if tt.findOneErr {
				mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &room).Return(assert.AnError)
			} else {
				mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &room).Return(nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)
			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			_, err := mockSvcStruct.GetRoomByID(tt.request, mongoPkgMock)
			if (err != nil) != tt.returnErr {
				t.Errorf("GetRoomByID() [%s] error = %v, initErr %v, returnErr %v", tt.name, err, tt.initErr, tt.returnErr)
			}
			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestJoinRoom(t *testing.T) {
	tests := []struct {
		name         string
		initErr      bool
		request      string
		updateOneErr bool
		returnErr    bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4567", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4567", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"updateone_error", false, "64a7b2f4e13e4c3f9c8b4567", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			if tt.updateOneErr {
				mongoCollectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, assert.AnError)
			} else {
				mongoCollectionMock.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)

			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			err := mockSvcStruct.JoinRoom(tt.request, 1, mongoPkgMock)
			if (err != nil) != tt.returnErr {
				t.Errorf("JoinRoom() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}
			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestGetRoomsForAll(t *testing.T) {
	tests := []struct {
		name       string
		initErr    bool
		request    string
		findOneErr bool
		decodeErr  bool
		returnErr  bool
	}{
		{"success_all", false, "all", false, false, false},
		{"success_joined", false, "joined", false, false, false},
		{"error", true, "all", false, false, true},
		{"invalid_target", false, "invalid_target", false, false, true},
		{"findone_error", false, "all", true, false, true},
		{"decode_error", false, "all", false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			var room model.Room
			mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)
			mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
			mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
			if tt.decodeErr {
				mongoCursorMock.On("Decode", &room).Return(assert.AnError)
			} else {
				mongoCursorMock.On("Decode", &room).Return(nil)
			}
			mongoCursorMock.On("Close", mock.Anything).Return(nil)

			var filter bson.M
			if tt.request == "all" {
				filter = bson.M{
					"$or": []bson.M{
						{"isprivate": false},
						{"members": 1},
					},
				}
			} else {
				filter = bson.M{"members": 1} // 参加済みのものだけ
			}

			if tt.findOneErr {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, assert.AnError)
			} else {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", "rooms").Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}

			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)

			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			_, err := mockSvcStruct.GetRooms(1, tt.request, mongoPkgMock)

			if (err != nil) != tt.returnErr {
				t.Errorf("GetRooms() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestPostChatMessage(t *testing.T) {
	tests := []struct {
		name          string
		initErr       bool
		requestRoomId string
		insertOneErr  bool
		returnErr     bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4567", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4567", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"insertone_error", false, "64a7b2f4e13e4c3f9c8b4567", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			if tt.insertOneErr {
				mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("", assert.AnError)
			} else {
				mongoCollectionMock.On("InsertOne", mock.Anything, mock.Anything).Return("mocked_id", nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", model.ChatMessageCollectionName).Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)

			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			err := mockSvcStruct.PostChatMessage(
				tt.requestRoomId,
				1,
				"Hello, World!",
				mongoPkgMock,
			)

			if (err != nil) != tt.returnErr {
				t.Errorf("PostChatMessage() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestGetChatMessages(t *testing.T) {
	tests := []struct {
		name      string
		initErr   bool
		findErr   bool
		decodeErr bool
		returnErr bool
	}{
		{"success", false, false, false, false},
		{"error", true, false, false, true},
		{"find_error", false, true, false, true},
		{"decode_error", false, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			var chatMessage model.ChatMessage
			mongoCursorMock := new(mock_mongo_pkg.MongoCursorMock)
			mongoCursorMock.On("Next", mock.Anything).Return(true).Once()
			mongoCursorMock.On("Next", mock.Anything).Return(false).Once()
			if tt.decodeErr {
				mongoCursorMock.On("Decode", &chatMessage).Return(assert.AnError)
			} else {
				mongoCursorMock.On("Decode", &chatMessage).Return(nil)
			}
			mongoCursorMock.On("Close", mock.Anything).Return(nil)

			filter := bson.M{"roomid": "64a7b2f4e13e4c3f9c8b4567"}

			if tt.findErr {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, assert.AnError)
			} else {
				mongoCollectionMock.On("Find", mock.Anything, filter).Return(mongoCursorMock, nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", model.ChatMessageCollectionName).Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}

			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)
			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
			_, err := mockSvcStruct.GetChatMessages("64a7b2f4e13e4c3f9c8b4567", mongoPkgMock)
			if (err != nil) != tt.returnErr {
				t.Errorf("GetChatMessages() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestReadChatMessages(t *testing.T) {
	tests := []struct {
		name          string
		initErr       bool
		requestChatId string
		updateErr     bool
		returnErr     bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4567", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4567", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"update_error", false, "64a7b2f4e13e4c3f9c8b4567", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			if tt.updateErr {
				mongoCollectionMock.On("UpdateMany", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, assert.AnError)
			} else {
				mongoCollectionMock.On("UpdateMany", mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", model.ChatMessageCollectionName).Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}

			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)
			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)
			err := mockSvcStruct.ReadChatMessages(
				"64a7b2f4e13e4c3f9c8b4567",
				[]string{
					tt.requestChatId,
					"64a7b2f4e13e4c3f9c8b4569",
				},
				1,
				mongoPkgMock,
			)

			if (err != nil) != tt.returnErr {
				t.Errorf("ReadChatMessages() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestGetChatMessageByID(t *testing.T) {
	tests := []struct {
		name             string
		initErr          bool
		requestMessageId string
		findOneErr       bool
		returnErr        bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4568", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4568", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"findone_error", false, "64a7b2f4e13e4c3f9c8b4568", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			var chatMessage model.ChatMessage
			if tt.findOneErr {
				mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &chatMessage).Return(assert.AnError)
			} else {
				mongoCollectionMock.On("FindOne", mock.Anything, mock.Anything, &chatMessage).Return(nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", model.ChatMessageCollectionName).Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)

			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			_, err := mockSvcStruct.GetChatMessageByID("64a7b2f4e13e4c3f9c8b4567", tt.requestMessageId, mongoPkgMock)

			if (err != nil) != tt.returnErr {
				t.Errorf("GetChatMessageByID() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteChatMessageByID(t *testing.T) {
	tests := []struct {
		name             string
		initErr          bool
		requestMessageId string
		deleteErr        bool
		returnErr        bool
	}{
		{"success", false, "64a7b2f4e13e4c3f9c8b4568", false, false},
		{"error", true, "64a7b2f4e13e4c3f9c8b4568", false, true},
		{"invalid_id", false, "invalid_object_id", false, true},
		{"delete_error", false, "64a7b2f4e13e4c3f9c8b4568", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoCollectionMock := new(mock_mongo_pkg.MongoCollectionMock)
			if tt.deleteErr {
				mongoCollectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, assert.AnError)
			} else {
				mongoCollectionMock.On("DeleteOne", mock.Anything, mock.Anything).Return(&mongo.DeleteResult{}, nil)
			}
			mongoDatabaseMock := new(mock_mongo_pkg.MongoDatabaseMock)
			mongoDatabaseMock.On("Collection", model.ChatMessageCollectionName).Return(mongoCollectionMock)

			mongoPkgStruct := &mongo_pkg.MongoPkgStruct{
				Ctx:    context.TODO(),
				Db:     mongoDatabaseMock,
				Cancel: func() {},
			}
			mongoPkgMock := setupInitMock(tt.initErr, "chatapp", mongoPkgStruct)

			mockSvcStruct := NewMongoSvc(mongoDatabaseMock)

			err := mockSvcStruct.DeleteChatMessage("64a7b2f4e13e4c3f9c8b4567", tt.requestMessageId, mongoPkgMock)

			if (err != nil) != tt.returnErr {
				t.Errorf("DeleteChatMessageByID() [%s] error = %v, initErr %v", tt.name, err, tt.initErr)
			}

			if m, ok := mongoPkgMock.(interface{ AssertExpectations(*testing.T) }); ok {
				m.AssertExpectations(t)
			}
		})
	}
}
