package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/pkg/mongo_pkg"
	"microservices/chat/tests/mocks/svc/mock_chat_svc"
	"microservices/chat/tests/mocks/svc/mock_mongo_svc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MongoPkgMock struct{}

func (m *MongoPkgMock) NewMongoConnect(dbName string) (*mongo_pkg.MongoPkgStruct, error) {
	return &mongo_pkg.MongoPkgStruct{}, nil
}

func TestCreateRoomHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock request
	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	mongoMockSvc.On("CreateRoom", mock.MatchedBy(func(r model.Room) bool {
		return r.Name == "TestRoom"
	}), mongoMockPkg).Return("mocked_room_id", nil)

	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader("name=TestRoom")
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoomHandler_InvalidRequest(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock request
	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)

	body := strings.NewReader("") // Missing 'name' field
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRoomHandler_DBError(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock request
	mongoMockPkg := &MongoPkgMock{}
	mongoMockSvc := new(mock_mongo_svc.MongoSvcMockWithErrorMock)
	chatMockSvc := new(mock_chat_svc.ChatSvcMock)
	mongoMockSvc.On("CreateRoom", mock.Anything, mongoMockPkg).Return(nil, assert.AnError)

	body := strings.NewReader("name=TestRoom")
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mongoMockSvc, mongoMockPkg, chatMockSvc)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
