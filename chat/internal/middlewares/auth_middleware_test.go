package middlewares

import (
	"crypto/rand"
	"crypto/rsa"
	"microservices/chat/internal/svc/jwtinfo_svc"
	"microservices/chat/tests/test_funcs"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestExtractBearerTokenFail(t *testing.T) {
	m := NewAuthMiddleware()

	r := gin.New()
	r.Use(m.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "not set jwt token")
}

func TestAuthMiddleware_Success(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "jwt_secret_key",
	}

	test_funcs.WithEnvMap(envs, t, func() {
		m := NewAuthMiddleware()

		r := gin.New()
		r.Use(m.Handler())
		jwt, err := test_funcs.CreateMockJwtToken(
			1,
			"test@example.com",
			time.Now().Add(1*time.Hour),
			[]byte("jwt_secret_key"),
		)
		if err != nil {
			t.Fatalf("failed to create mock JWT token: %v", err)
		}
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwt)
		w := httptest.NewRecorder()
		r.GET("/test", func(c *gin.Context) {
			userID := c.Request.Context().Value(jwtinfo_svc.UserIDKey)
			assert.Equal(t, userID, 1)
			email := c.Request.Context().Value(jwtinfo_svc.EmailKey)
			assert.Equal(t, email, "test@example.com")
			c.JSON(200, gin.H{"userID": 1, "email": "test@example.com"})
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "jwt_secret_key",
	}

	test_funcs.WithEnvMap(envs, t, func() {
		m := NewAuthMiddleware()

		r := gin.New()
		r.Use(m.Handler())

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")

		w := httptest.NewRecorder()
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "invalid jwt token")
	})
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "jwt_secret_key",
	}

	test_funcs.WithEnvMap(envs, t, func() {
		m := NewAuthMiddleware()

		r := gin.New()
		r.Use(m.Handler())

		jwt, err := test_funcs.CreateMockJwtToken(
			1,
			"test@example.com",
			time.Now().Add(-1*time.Hour),
			[]byte("jwt_secret_key"),
		)
		if err != nil {
			t.Fatalf("failed to create mock JWT token: %v", err)
		}
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		w := httptest.NewRecorder()
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "token is expired")
	})
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "wrong_secret_key",
	}

	test_funcs.WithEnvMap(envs, t, func() {
		m := NewAuthMiddleware()

		r := gin.New()
		r.Use(m.Handler())

		jwt, err := test_funcs.CreateMockJwtToken(
			1,
			"test@example.com",
			time.Now().Add(1*time.Hour),
			[]byte("jwt_secret_key"),
		)
		if err != nil {
			t.Fatalf("failed to create mock JWT token: %v", err)
		}
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		w := httptest.NewRecorder()
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "invalid jwt token")
	})
}

func TestSigningMethodHMACFail(t *testing.T) {
	envs := test_funcs.Envs{
		"JWT_SECRET": "jwt_secret_key",
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	claims := jwt.MapClaims{
		"sub":   1,
		"email": "wrong@sig.example.com",
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	// RS256 で署名！
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("failed to create mock JWT token: %v", err)
	}

	test_funcs.WithEnvMap(envs, t, func() {
		m := NewAuthMiddleware()

		r := gin.New()
		r.Use(m.Handler())
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		w := httptest.NewRecorder()
		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "success"})
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "unexpected signing method")
	})
}
