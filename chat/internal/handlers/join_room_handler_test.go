package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/tests/mocks/svc/mock_chat_svc"
	"microservices/chat/tests/mocks/svc/mock_mongo_svc"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJoinRoomHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)

	mongoMockSvc.On("JoinRoom", "valid_room_id", int(12345), mongoMockPkg).Return(nil)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Joined room successfully")
}

func TestJoinRoomHandlerGetRoomByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, assert.AnError)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestJoinRoomHandlerJoinRoomError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("GetRoomByID", "valid_room_id", mongoMockPkg).Return(model.Room{}, nil)
	mongoMockSvc.On("JoinRoom", "valid_room_id", int(12345), mongoMockPkg).Return(assert.AnError)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to join room")
}

func TestJoinRoomHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader("") // Missing 'room_id' field
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}
