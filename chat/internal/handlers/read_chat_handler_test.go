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
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReadChatMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("ReadChatMessages", "valid_room_id", []string{"chat1", "chat2"}, int(12345), mongoMockPkg).Return(nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: true, IsOwner: false})

	body := strings.NewReader(`{
		"room_id": "valid_room_id",
		"chat_id_list": ["chat1", "chat2"]
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Chat messages marked as read")
}

func TestReadChatMessagesGetRoomByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, assert.AnError)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader(`{
		"room_id": "valid_room_id",
		"chat_id_list": ["chat1", "chat2"]
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestReadChatMessagesGetRoomInfoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: false, IsOwner: false})

	body := strings.NewReader(`{
		"room_id": "valid_room_id",
		"chat_id_list": ["chat1", "chat2"]
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Access denied")
}

func TestReadChatMessagesZeroLen(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: true, IsOwner: false})

	body := strings.NewReader(`{
		"room_id": "valid_room_id",
		"chat_id_list": []
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "chat_id_list is required")
}

func TestReadChatMessagesReadChatMessagesError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("ReadChatMessages", "valid_room_id", []string{"chat1", "chat2"}, int(12345), mongoMockPkg).Return(assert.AnError)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	chatMockSvc.On("GetRoomInfo", model.Room{}, int(12345)).Return(chat_svc.Room{IsMember: true, IsOwner: false})

	body := strings.NewReader(`{
		"room_id": "valid_room_id",
		"chat_id_list": ["chat1", "chat2"]
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to read chat message")
}

func TestReadChatMessagesInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader(`{
		"room_id": "",
		"chat_id_list": []
	}`)
	req := httptest.NewRequest("POST", "/read_chat", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.ReadChatMessages(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}
