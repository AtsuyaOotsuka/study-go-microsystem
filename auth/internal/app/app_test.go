package app_test

import (
	"microservices/auth/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	// 必要なら  AutoMigrate(&models.User{})
	return db
}

func TestNewAppAndInitRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// App生成とルート初期化
	db := newTestDB(t) // ← 本物の *gorm.DB
	sqlDB, _ := db.DB()
	a, cleanup, err := app.NewApp(db, sqlDB) // ← これで db.DB() もOK
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	a.InitRoutes(r)

	// テスト用リクエスト
	req := httptest.NewRequest("GET", "/health", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 結果検証
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}
