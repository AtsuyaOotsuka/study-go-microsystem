package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"microservices/auth/pkg/csrf_pkg"
	"microservices/auth/tests/test_funcs"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testServer *http.Server
var baseURL string
var db *gorm.DB
var sqlDB *sql.DB
var dbRecords []test_funcs.DbRecords

func TestMain(m *testing.M) {
	godotenv.Load(".env.test")
	gin.SetMode(gin.TestMode)
	db, sqlDB = SetupDB()

	var err error

	dbRecords, err = test_funcs.DbCleanup(sqlDB)
	assert.NoError(&testing.T{}, err)

	r, cleanup := SetupRouter(db, sqlDB)
	defer cleanup()

	testServer = &http.Server{
		Addr:    ":8880",
		Handler: r,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
	baseURL = "http://localhost:8880"

	fmt.Println("Test server started at", baseURL)

	// 全テスト実行
	exitCode := m.Run()

	// サーバをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = testServer.Shutdown(ctx)

	os.Exit(exitCode)
}

func createCsrf() string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := test_funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func request(method string, url string, body io.Reader, t *testing.T) (*http.Response, func() error) {
	csrf := createCsrf()

	client := &http.Client{}
	requestUrl := baseURL + url
	fmt.Println("Request URL:", requestUrl)
	req, err := http.NewRequest(method, requestUrl, body)
	if method != "GET" && err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	assert.NoError(t, err)
	if method != "GET" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, resp.Body.Close
}

func TestHealth(t *testing.T) {
	resp, close := request("GET", "/health", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	assert.Contains(t, bodyString, "healthy")
}

func TestLogin(t *testing.T) {

	id := dbRecords[0].Data[0]["id"].(int64)
	email := dbRecords[0].Data[0]["email"].(string)
	password := dbRecords[0].Data[0]["password"].(string)
	refreshToken := dbRecords[0].Data[0]["refresh_token"].(string)

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	resp, close := request("POST", "/auth/login", io.NopCloser(strings.NewReader(body)), t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)
	fmt.Println("access_token", respData["access_token"])

	jwt := respData["access_token"].(string)
	jwtInfo, err := test_funcs.JwtConvert(jwt)
	assert.NoError(t, err)

	assert.Equal(t, int(id), jwtInfo.UserID)
	assert.Equal(t, email, jwtInfo.Email)
	assert.Equal(t, refreshToken, respData["refresh_token"])
}

func TestRefresh(t *testing.T) {
	id := dbRecords[0].Data[0]["id"].(int64)
	email := dbRecords[0].Data[0]["email"].(string)
	refreshToken := dbRecords[0].Data[0]["refresh_token"].(string)

	body := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	resp, close := request("POST", "/auth/refresh", io.NopCloser(strings.NewReader(body)), t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)
	fmt.Println("access_token", respData["access_token"])

	jwt := respData["access_token"].(string)
	jwtInfo, err := test_funcs.JwtConvert(jwt)
	assert.NoError(t, err)

	assert.Equal(t, int(id), jwtInfo.UserID)
	assert.Equal(t, email, jwtInfo.Email)
	assert.Equal(t, refreshToken, respData["refresh_token"]) // リフレッシュトークンの更新は現状はしていないので、同じ値
}

func TestRegister(t *testing.T) {
	body := `{"name":"New User","email":"newuser@example.com","password":"password123"}`
	resp, close := request("POST", "/register", io.NopCloser(strings.NewReader(body)), t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	insertResult := test_funcs.ExistsRecord(sqlDB, "users", map[string]interface{}{
		"email": "newuser@example.com",
	})
	assert.True(t, insertResult)
}

func TestCsrfGet(t *testing.T) {
	resp, close := request("GET", "/csrf/get", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData["csrf_token"])

	cookies := resp.Cookies()
	var csrfCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "csrf_token" {
			csrfCookie = cookie
			break
		}
	}
	assert.NotNil(t, csrfCookie)
	assert.NotEmpty(t, csrfCookie.Value)

	// Cookieから出したため、URLデコードを実施
	decodeToken, err := url.QueryUnescape(csrfCookie.Value)
	assert.NoError(t, err)

	// クッキーの値とレスポンスの値が同じであることを確認
	assert.Equal(t, decodeToken, respData["csrf_token"])

	csrfStruct := csrf_pkg.CsrfPkgStruct{}

	// CSRFトークンの検証
	err = csrfStruct.ValidateCSRFCookieToken(
		decodeToken,
		os.Getenv("CSRF_TOKEN"),
		time.Now().Unix(),
	)
	assert.NoError(t, err)
}
