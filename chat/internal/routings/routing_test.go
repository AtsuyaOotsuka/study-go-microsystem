package routings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockHandlers struct{}

func (m *MockHandlers) CreateRoomHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}
func (m *MockHandlers) HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}
func (m *MockHandlers) JoinRoomHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}
func (m *MockHandlers) RoomListHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}
func (m *MockHandlers) PostChatMessageHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

type MockMiddleware struct{}

func (m *MockMiddleware) CSRFMW(c *gin.Context) {
	// Mock CSRF middleware logic
	c.Next()
}

func (m *MockMiddleware) AuthMW(c *gin.Context) {
	// Mock Auth middleware logic
	c.Next()
}

func TestRouting(t *testing.T) {
	expected := map[string]string{
		"/room_create": "POST",
	}

	r := gin.Default()
	mwMock := &MockMiddleware{}
	handlersMock := &MockHandlers{}
	Routing(r, mwMock.CSRFMW, mwMock.AuthMW, handlersMock)

	for path, method := range expected {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, `{"status": "success"}`, w.Body.String())
		})
	}
}
