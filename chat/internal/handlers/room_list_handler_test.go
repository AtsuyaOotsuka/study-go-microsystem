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

func TestRoomListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRooms", int(12345), "all", mongoMockPkg).Return([]model.Room{}, nil)
	chatMockSvc.On("ConvertRoomList", []model.Room{}, int(12345)).Return([]chat_svc.Room{})

	req := httptest.NewRequest("GET", "/rooms?target=all", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req
	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.RoomListHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"rooms":[]`)
}

func TestRoomListHandler_GetRoomsError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	mongoMockSvc.On("GetRooms", int(12345), "all", mongoMockPkg).Return([]model.Room{}, assert.AnError)

	req := httptest.NewRequest("GET", "/rooms?target=all", nil)
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)
	c.Request = req
	handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	handler.RoomListHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `Failed to get rooms`)
}
