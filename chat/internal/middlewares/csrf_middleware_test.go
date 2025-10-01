package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCSRFVerifier struct {
	mock.Mock
}

func (m *MockCSRFVerifier) Verify(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestCsrfHandler(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "valid_token").Return(nil)

	middleware := NewCSRFMiddleware(mockVerifier)
	handler := middleware.Handler()

	r := gin.New()
	r.Use(handler)
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not set csrf token")
}

func TestCSRFMiddleware_InvalidToken(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "invalid_token").Return(assert.AnError)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "invalid_token")
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
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	mockVerifier.AssertNotCalled(t, "Verify", mock.Anything, mock.Anything)
}

func TestCSRFMiddleware_forCookie(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "cookie_token").Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "cookie_token"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	mockVerifier.AssertCalled(t, "Verify", mock.Anything, "cookie_token")
}

func TestCSRFMiddlewareSuccess(t *testing.T) {
	mockVerifier := new(MockCSRFVerifier)
	mockVerifier.On("Verify", mock.Anything, "valid_token").Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockVerifier).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "valid_token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	mockVerifier.AssertCalled(t, "Verify", mock.Anything, "valid_token")
}
