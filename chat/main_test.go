package main

import (
	"context"
	"fmt"
	"io"
	"microservices/chat/tests/test_funcs"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testServer *http.Server
var baseURL string

func TestMain(m *testing.M) {

	godotenv.Load(".env.test")

	// 起動前のセットアップ
	gin.SetMode(gin.TestMode)
	r := SetupRouter()
	testServer = &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// サーバが立ち上がるまで待つ（またはヘルスチェックする）
	time.Sleep(200 * time.Millisecond)
	baseURL = "http://localhost:8081"

	fmt.Println("Test server started at", baseURL)

	// 全テスト実行
	exitCode := m.Run()

	// サーバをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = testServer.Shutdown(ctx)

	os.Exit(exitCode)
}

func createJwt(t *testing.T) string {
	jwt_secret := os.Getenv("JWT_SECRET")

	jwt, err := test_funcs.CreateMockJwtToken(
		1,
		"testuser@example.com",
		time.Now().Add(1*time.Hour),
		[]byte(jwt_secret),
	)

	assert.NoError(t, err)

	return jwt
}

func createCsrf(t *testing.T) string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := test_funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func request(method string, url string, body io.Reader, t *testing.T) (*http.Response, error, func() error) {
	jwt := createJwt(t)
	csrf := createCsrf(t)

	client := &http.Client{}
	req, err := http.NewRequest(method, baseURL+url, body)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("X-CSRF-Token", csrf)
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, nil, resp.Body.Close
}

func TestHealth(t *testing.T) {
	resp, err, close := request("GET", "/health", nil, t)
	defer close()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "healthy")
}
