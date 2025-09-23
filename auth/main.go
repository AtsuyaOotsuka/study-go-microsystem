package main

import (
	"log"
	"microservices/auth/internal/app"
	"microservices/auth/pkg/db_pkg"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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
	r := gin.New()
	r.Use(gin.Recovery())

	db_pkg := db_pkg.NewDBConnect()
	db, _ := db_pkg.ConnectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	app, cleanup, _ := app.NewApp(db, sqlDB)
	defer cleanup()

	app.InitRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
