package routings

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCSRFHandler struct{}

func (m *MockCSRFHandler) CsrfGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func TestCsrfRouting(t *testing.T) {
	expected := map[string]string{
		"/csrf/get": "GET",
	}

	r := gin.Default()
	CsrfRouting(r, &MockCSRFHandler{})

	for path, method := range expected {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, `{"status": "OK"}`, w.Body.String())
		})
	}
}
