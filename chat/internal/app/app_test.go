package app

import (
	"microservices/chat/tests/test_funcs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitRoutes(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "jwt_secret_key",
		"CSRF_TOKEN": "csrf_token",
	}

	test_funcs.WithEnvMap(envs, t, func() {

		csrf := test_funcs.GenerateCSRFCookieToken(
			"csrf_token",
			time.Now().Add(1*time.Hour).Unix(),
		)
		jwt, err := test_funcs.CreateMockJwtToken(
			1,
			"user@example.com",
			time.Now().Add(1*time.Hour),
			[]byte("jwt_secret_key"),
		)
		if err != nil {
			t.Fatalf("failed to create mock JWT token: %v", err)
		}

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		app := NewApp()
		app.InitRoutes(r)

		req := httptest.NewRequest("GET", "/health", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CSRF-Token", csrf)
		req.Header.Set("Authorization", "Bearer "+jwt)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		// 結果検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "healthy")
	})
}
