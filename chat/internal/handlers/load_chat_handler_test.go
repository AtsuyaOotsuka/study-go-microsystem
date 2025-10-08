package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/chat_svc"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/tests/mocks/svc/mock_chat_svc"
	"microservices/chat/tests/mocks/svc/mock_mongo_svc"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoadChatHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("GetChatMessages", "valid_room_id", mongoMockPkg).Return([]model.ChatMessage{}, nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: true, IsOwner: false})

	req := httptest.NewRequest("GET", "/load_chat/valid_room_id", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	c.Params = append(c.Params, gin.Param{Key: "room_id", Value: "valid_room_id"})
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.LoadChatHandlers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "messages")
}

func TestLoadChatHandlersGetRoomByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, assert.AnError)

	req := httptest.NewRequest("GET", "/load_chat/valid_room_id", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	c.Params = append(c.Params, gin.Param{Key: "room_id", Value: "valid_room_id"})
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.LoadChatHandlers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestLoadChatHandlersGetRoomInfoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: false, IsOwner: false})

	req := httptest.NewRequest("GET", "/load_chat/valid_room_id", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	c.Params = append(c.Params, gin.Param{Key: "room_id", Value: "valid_room_id"})
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.LoadChatHandlers(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

func TestLoadChatHandlersGetChatMessagesError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("GetChatMessages", "valid_room_id", mongoMockPkg).Return([]model.ChatMessage{}, assert.AnError)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: true, IsOwner: false})

	req := httptest.NewRequest("GET", "/load_chat/valid_room_id", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	c.Params = append(c.Params, gin.Param{Key: "room_id", Value: "valid_room_id"})
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.LoadChatHandlers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get chat messages")
}
