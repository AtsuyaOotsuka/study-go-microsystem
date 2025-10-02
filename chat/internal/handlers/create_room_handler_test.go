package handlers

import (
	"context"
	"microservices/chat/internal/model"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/pkg/mongo_pkg"
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
	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	mockSvc.On("CreateRoom", mock.MatchedBy(func(r model.Room) bool {
		return r.Name == "TestRoom"
	}), mockPkg).Return("mocked_room_id", nil)
	c.Set("mongo", mockSvc)

	body := strings.NewReader("name=TestRoom")
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mockSvc, mockPkg)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateRoomHandler_InvalidRequest(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock request
	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMock)
	c.Set("mongo", mockSvc)

	body := strings.NewReader("") // Missing 'name' field
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mockSvc, mockPkg)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateRoomHandler_DBError(t *testing.T) {

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Mock request
	mockPkg := &MongoPkgMock{}
	mockSvc := new(mock_mongo_svc.MongoSvcMockWithErrorMock)
	mockSvc.On("CreateRoom", mock.Anything, mockPkg).Return(nil, assert.AnError)
	c.Set("mongo", mockSvc)

	body := strings.NewReader("name=TestRoom")
	req := httptest.NewRequest("POST", "/rooms", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.WithValue(req.Context(), jwtinfo_svc.UserIDKey, 12345)
	ctx = context.WithValue(ctx, jwtinfo_svc.EmailKey, "test@example.com")
	req = req.WithContext(ctx)

	c.Request = req

	Handler := NewHandlers(mockSvc, mockPkg)
	Handler.CreateRoomHandler(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
