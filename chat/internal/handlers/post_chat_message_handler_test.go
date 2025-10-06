package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/tests/mocks/svc/mock_mongo_svc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostChatMessageHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("PostChatMessage", "valid_room_id", int(12345), "Hello, World!", mongoMockPkg).Return(nil)

	body := strings.NewReader(`{"room_id":"valid_room_id","message":"Hello, World!"}`)
	req := httptest.NewRequest("POST", "/post_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, nil)
	handler.PostChatMessageHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Chat posted successfully")
}

func TestPostChatMessageHandlerGetRoomByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, assert.AnError)

	body := strings.NewReader(`{"room_id":"valid_room_id","message":"Hello, World!"}`)
	req := httptest.NewRequest("POST", "/post_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, nil)
	handler.PostChatMessageHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestPostChatMessageHandlerPostChatMessageError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("PostChatMessage", "valid_room_id", int(12345), "Hello, World!", mongoMockPkg).Return(assert.AnError)

	body := strings.NewReader(`{"room_id":"valid_room_id","message":"Hello, World!"}`)
	req := httptest.NewRequest("POST", "/post_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, nil)
	handler.PostChatMessageHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to post chat")
}

func TestPostChatMessageHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)

	body := strings.NewReader(`{"room_id":"","message":""}`)
	req := httptest.NewRequest("POST", "/post_chat_message", body)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, nil)
	handler.PostChatMessageHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}
