package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
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

	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	mockSvc.On("GetRoomByID", "valid_room_id", mockPkg).Return(model.Room{}, nil)

	mockSvc.On("JoinRoom", "valid_room_id", int(12345), mockPkg).Return(nil)
	c.Set("mongo", mockSvc)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mockSvc, mockPkg)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Joined room successfully")
}

func TestJoinRoomHandlerGetRoomByIDError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	mockSvc.On("GetRoomByID", "valid_room_id", mockPkg).Return(model.Room{}, assert.AnError)

	c.Set("mongo", mockSvc)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mockSvc, mockPkg)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to get room")
}

func TestJoinRoomHandlerJoinRoomError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	mockSvc.On("GetRoomByID", "valid_room_id", mockPkg).Return(model.Room{}, nil)
	mockSvc.On("JoinRoom", "valid_room_id", int(12345), mockPkg).Return(assert.AnError)

	c.Set("mongo", mockSvc)

	body := strings.NewReader("room_id=valid_room_id")
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mockSvc, mockPkg)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to join room")
}

func TestJoinRoomHandlerInvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	c.Set("mongo", mockSvc)

	body := strings.NewReader("") // Missing 'room_id' field
	req := httptest.NewRequest("POST", "/join_room", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req

	handler := NewHandlers(mockSvc, mockPkg)
	handler.JoinRoomHandler(c)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")
}
