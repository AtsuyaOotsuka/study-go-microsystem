package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CSRFVerifier interface {
	Verify(ctx context.Context, token string) error
}

type CSRFMiddleware struct{ Verifier CSRFVerifier }

func NewCSRFMiddleware(v CSRFVerifier) *CSRFMiddleware {
	return &CSRFMiddleware{Verifier: v}
}

func (m *CSRFMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Next()
			return
		}
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("_token")
		}
		if token == "" {
			cookie, err := c.Cookie("csrf_token")
			if err == nil {
				token = cookie
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "not set csrf token"})
			return
		}
		if err := m.Verifier.Verify(c.Request.Context(), token); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			return
		}
		c.Next()
	}
}
