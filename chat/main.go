package main

import (
	"log"
	"microservices/chat/internal/app"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	app := app.NewApp()
	app.InitRoutes(r)

	return r
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
	r := SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
