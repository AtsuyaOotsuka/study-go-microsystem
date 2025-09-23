package routings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockAuthHandler struct{}

func (m *MockAuthHandler) HandleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (m *MockAuthHandler) HandleRefresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func TestAuthRouting(t *testing.T) {
	expected := map[string]string{
		"/auth/login":   "POST",
		"/auth/refresh": "POST",
	}

	r := gin.Default()
	AuthRouting(r, &MockAuthHandler{}, func(c *gin.Context) {
		c.Next()
	})

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
