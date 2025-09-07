package handlers

import (
	"fmt"
	"microservices/auth/internal/clock_svc"
	"microservices/auth/tests/mocks/global_mock"
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

	gdb, mock := global_mock.NewGormWithMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	body := strings.NewReader("name=Test+User&email=test%40example.com&password=password123")
	req := httptest.NewRequest("POST", "/auth/register", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = req

	handler := NewRegisterHandler(
		gdb, &jwt.JwtServiceMockStruct{}, clock_svc.RealClockStruct{},
	)
	handler.HandleRegister(c)

	fmt.Println("📣 recorder body:", w.Body.String()) // ←追加

	assert.Equal(t, http.StatusOK, w.Code)
}
