package main

import (
	"database/sql"
	"log"
	"microservices/auth/internal/app"
	"microservices/auth/pkg/db_pkg"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, sqlDB *sql.DB) (*gin.Engine, func()) {
	r := gin.New()
	r.Use(gin.Recovery())

	app, cleanup, _ := app.NewApp(db, sqlDB)

	app.InitRoutes(r)

	return r, cleanup
}

func SetupDB() (*gorm.DB, *sql.DB) {
	db_pkg := db_pkg.NewDBConnect()
	db, _ := db_pkg.ConnectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	return db, sqlDB
}

func main() {
	// .envを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ .envファイル読み込み失敗（デフォルトのdebugモードで起動）")
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode // fallback
	}
	gin.SetMode(mode)

	db, sqlDB := SetupDB()
	r, cleanup := SetupRouter(db, sqlDB)
	defer cleanup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
