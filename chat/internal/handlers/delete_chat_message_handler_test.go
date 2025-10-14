package handlers

import (
	"context"
	"microservices/chat/internal/model"
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

func TestDeleteChatMessageHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("GetChatMessageByID", "valid_room_id", "valid_message_id", mongoMockPkg).Return(model.ChatMessage{UserID: 12345}, nil)
	mongoMockSvc.On("DeleteChatMessage", "valid_room_id", "valid_message_id", mongoMockPkg).Return(nil)

	body := strings.NewReader(`{"room_id":"valid_room_id","message_id":"valid_message_id"}`)
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Message deleted successfully")
}

func TestDeleteChatMessageHandlerFaildGetRoom(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, assert.AnError)

	body := strings.NewReader(`{"room_id":"valid_room_id","message_id":"valid_message_id"}`)
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestDeleteChatMessageHandlerFaildGetMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("GetChatMessageByID", "valid_room_id", "valid_message_id", mongoMockPkg).Return(model.ChatMessage{}, assert.AnError)

	body := strings.NewReader(`{"room_id":"valid_room_id","message_id":"valid_message_id"}`)
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get message")
}

func TestDeleteChatMessageHandlerForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{OwnerID: 99999}, nil)
	mongoMockSvc.On("GetChatMessageByID", "valid_room_id", "valid_message_id", mongoMockPkg).Return(model.ChatMessage{UserID: 12345}, nil)

	body := strings.NewReader(`{"room_id":"valid_room_id","message_id":"valid_message_id"}`)
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 54321) // Different user
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "You can only delete your own messages")
}

func TestDeleteChatMessageHandlerFailedDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("GetChatMessageByID", "valid_room_id", "valid_message_id", mongoMockPkg).Return(model.ChatMessage{UserID: 12345}, nil)
	mongoMockSvc.On("DeleteChatMessage", "valid_room_id", "valid_message_id", mongoMockPkg).Return(assert.AnError)

	body := strings.NewReader(`{"room_id":"valid_room_id","message_id":"valid_message_id"}`)
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete message")
}

func TestDeleteChatMessageHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader(`{"room_id":"","message_id":"valid_message_id"}`) // Missing room_id
	req := httptest.NewRequest("DELETE", "/delete_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "user@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.DeleteChatMessageHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}
