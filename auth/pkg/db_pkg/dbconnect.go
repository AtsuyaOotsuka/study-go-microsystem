package db_pkg

import (
	"fmt"
	"microservices/auth/internal/models"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConnectInterface interface {
	ConnectDB() (*gorm.DB, error)
}

type DBConnectStruct struct{}

func NewDBConnect() *DBConnectStruct {
	return &DBConnectStruct{}
}

func (d *DBConnectStruct) ConnectDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	tz := os.Getenv("DB_TZ")

	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:               dbName,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset":      "utf8mb4",
			"parseTime":    "true",
			"loc":          tz,
			"timeout":      "5s",
			"readTimeout":  "5s",
			"writeTimeout": "5s",
		},
	}

	gcfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// SingularTable: true, // テーブル名を単数形にしたい場合は有効化
		},
		Logger:                                   logger.Default.LogMode(logger.Info), // LogLevel
		DisableForeignKeyConstraintWhenMigrating: true,                                // 生成時にFK作らない（好み）
	}

	db, err := gorm.Open(gormmysql.Open(cfg.FormatDSN()), gcfg)
	if err != nil {
		err = fmt.Errorf("DB接続失敗: %w", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("sqlDB取得失敗: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)                  // 同時接続上限
	sqlDB.SetMaxIdleConns(25)                  // アイドル保持
	sqlDB.SetConnMaxLifetime(60 * time.Minute) // 再作成サイクル
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // アイドルの寿命

	// Ping
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("DB Ping失敗: %w", err)
	}

	fmt.Println("DB接続成功")

	// マイグレーション (テーブル作成)
	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	// マイグレーション (テーブル作成)
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("マイグレーション失敗: %w", err)
	}
	return nil
}
