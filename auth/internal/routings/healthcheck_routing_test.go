package routings

import (
	"net/http"
	"testing"

	// ← ここをimportする！

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockHealthCheckHandler struct{}

func (m *MockHealthCheckHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func TestHealthCheckRouting(t *testing.T) {
	expected := map[string]string{
		"/health": "GET",
	}

	r := gin.Default()
	HealthCheckRouting(r, &MockHealthCheckHandler{})

	for path, method := range expected {
		found := false
		for _, route := range r.Routes() {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected route %s [%s] to be registered", path, method)
	}
}
