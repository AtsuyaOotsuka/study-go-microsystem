package routings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRegisterHandler struct{}

func (m *MockRegisterHandler) HandleRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func TestRegisterRouting(t *testing.T) {
	expected := map[string]string{
		"/register": "POST",
	}

	r := gin.Default()
	RegisterRouting(r, &MockRegisterHandler{}, func(c *gin.Context) {
		c.Next()
	})

	for path, method := range expected {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, `{"status": "registered"}`, w.Body.String())
		})
	}
}
