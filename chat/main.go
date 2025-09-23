package main

import (
	"log"
	"microservices/chat/internal/app"
	"microservices/chat/pkg/mongo_pkg"
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

	MongoPkgStruct, err := mongo_pkg.NewMongoConnect()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	defer MongoPkgStruct.Cancel()

	app := app.NewApp(MongoPkgStruct)
	app.InitRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
