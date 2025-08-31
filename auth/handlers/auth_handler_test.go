package handlers

import (
	"microservices/auth/tests/mocks/models_mock"
	"microservices/auth/tests/mocks/svc_internal/jwt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newGormWithMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return gdb, mock
}

func TestHandleLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockUser := models_mock.CreateUserMock()

	// sqlmock準備
	gdb, mock := newGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, mockUser.Password, mockUser.Email)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(mockUser.Email, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// POSTリクエストをセット
	body := strings.NewReader("email=test@example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleLogin(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mock.jwt.token")
}

func TestHandleLogin_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockUser := models_mock.CreateUserMock()

	// sqlmock準備
	gdb, mock := newGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, mockUser.Password, mockUser.Email)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(mockUser.Email, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// POSTリクエストをセット
	body := strings.NewReader("email=test@example.com&password=wrongpassword")
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid password")
}

func TestHandleLogin_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockUser := models_mock.CreateUserMock()

	// sqlmock準備
	gdb, mock := newGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(mockUser.Email, sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	// POSTリクエストをセット
	body := strings.NewReader("email=test@example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid email")
}

func TestHandleLogin_FailedValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// POSTリクエストをセット
	body := strings.NewReader("email=invalid-email&password=short")
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	gdb, _ := newGormWithMock(t)

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestHandleLogin_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockUser := models_mock.CreateUserMock()

	// sqlmock準備
	gdb, mock := newGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, mockUser.Password, mockUser.Email)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(mockUser.Email, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// POSTリクエストをセット
	body := strings.NewReader("email=test@example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceFailedMockStruct{}) // ★ モックDBを注入
	handler.HandleLogin(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create JWT")
}

func TestHandleRefresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// sqlmock準備
	gdb, mock := newGormWithMock(t)

	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs("mock_refresh_token", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "refresh_token"}).
			AddRow(1, "test@example.com", "mock_refresh_token"))

	// POSTリクエストをセット
	body := strings.NewReader("refresh_token=mock_refresh_token")
	req := httptest.NewRequest("POST", "/auth/refresh", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleRefresh(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mock.jwt.token")
}

func TestHandleRefresh_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// sqlmock準備
	gdb, mock := newGormWithMock(t)

	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs("invalid_refresh_token", sqlmock.AnyArg()).
		WillReturnError(gorm.ErrRecordNotFound)

	// POSTリクエストをセット
	body := strings.NewReader("refresh_token=invalid_refresh_token")
	req := httptest.NewRequest("POST", "/auth/refresh", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleRefresh(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid refresh token")
}

func TestHandleRefresh_NotRequestToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// POSTリクエストをセット
	body := strings.NewReader("not_refresh_token=xxxxxxxxx")
	req := httptest.NewRequest("POST", "/auth/refresh", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	gdb, _ := newGormWithMock(t)

	handler := NewAuthHandler(gdb, &jwt.JwtServiceMockStruct{}) // ★ モックDBを注入
	handler.HandleRefresh(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
}

func TestHandleRefresh_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	mockUser := models_mock.CreateUserMock()

	// sqlmock準備
	gdb, mock := newGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, mockUser.Password, mockUser.Email)

	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE refresh_token = \\?").
		WithArgs(mockUser.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(rows)

	// POSTリクエストをセット
	body := strings.NewReader("refresh_token=" + mockUser.RefreshToken)
	req := httptest.NewRequest("POST", "/auth/refresh", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewAuthHandler(gdb, &jwt.JwtServiceFailedMockStruct{}) // ★ モックDBを注入
	handler.HandleRefresh(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create JWT")
}
