package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"microservices/auth/tests/mocks/svc_internal/csrf"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewCSRFHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	csrf := &csrf.CsrfSvcMockStruct{}
	handler := NewCSRFHandler(csrf)
	handler.CsrfGet(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked_csrf_token")
}
