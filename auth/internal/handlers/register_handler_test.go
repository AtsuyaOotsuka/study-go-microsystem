package handlers

import (
	"microservices/auth/internal/svc/clock_svc"
	"microservices/auth/tests/mocks/global_mock"
	"microservices/auth/tests/mocks/pkg/encrypt_pkg_mock"
	"microservices/auth/tests/mocks/svc_internal/jwt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	defer cleanup()

	mock.ExpectCommit()
	body := strings.NewReader("name=Test+User&email=test%40example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewRegisterHandler(
		gdb,
		&jwt.JwtServiceMockStruct{},
		&encrypt_pkg_mock.EncryptPkgMockStruct{},
		clock_svc.RealClockStruct{},
	)
	handler.HandleRegister(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestShouldBindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := strings.NewReader("name=Test+User&email=invalid-email&password=short")
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewRegisterHandler(
		nil,
		&jwt.JwtServiceMockStruct{},
		&encrypt_pkg_mock.EncryptPkgMockStruct{},
		clock_svc.RealClockStruct{},
	)
	handler.HandleRegister(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestHashPasswordError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := strings.NewReader("name=Test+User&email=test%40example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewRegisterHandler(
		nil,
		&jwt.JwtServiceMockStruct{},
		&encrypt_pkg_mock.EncryptPkgMockErrorStruct{},
		clock_svc.RealClockStruct{},
	)
	handler.HandleRegister(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestFailCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// sqlmock準備
	gdb, mock, cleanup := global_mock.NewGormWithMockError(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()
	defer cleanup()

	body := strings.NewReader("name=Test+User&email=test%40example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewRegisterHandler(
		gdb,
		&jwt.JwtServiceMockStruct{},
		&encrypt_pkg_mock.EncryptPkgMockStruct{},
		clock_svc.RealClockStruct{},
	)
	handler.HandleRegister(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}
