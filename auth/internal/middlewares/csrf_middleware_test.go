package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCSRFVerifier struct {
	mock.Mock // ← これが大事！！
}

func (m *MockCSRFVerifier) Verify(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestCsrfHandler(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "valid_token").Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not set csrf token")
}

func TestCSRFMiddleware_InvalidToken(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.
		On("Verify", mock.Anything, "invalid-token").
		Return(fmt.Errorf("invalid"))

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid csrf token")
}

func TestCSRFMiddleware_forGET(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OK")
}

func TestCSRFMiddleware_forCookie(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "cookie_token").Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "cookie_token"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OK")
}

func TestCSRFMiddlewareSuccess(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "valid_token").Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "valid_token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "OK")
}
