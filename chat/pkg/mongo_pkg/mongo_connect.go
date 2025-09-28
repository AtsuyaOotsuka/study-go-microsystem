package mongo_pkg

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPkgStruct struct {
	Db     *mongo.Database
	Ctx    context.Context
	Cancel context.CancelFunc
}

func connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
	// タイムアウト付きのcontext
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)

	mongoURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		defer cancelFunc()
		return nil, nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		defer cancelFunc()
		return nil, nil, nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, ctx, cancelFunc, nil
}

func NewMongoConnect(database string) (*MongoPkgStruct, error) {
	client, ctx, cancelFunc, err := connect()
	if err != nil {
		return nil, err
	}

	m := &MongoPkgStruct{}
	m.Ctx = ctx
	m.Cancel = cancelFunc
	m.Db = client.Database(database)
	fmt.Println("Connected to MongoDB!")

	return m, nil
}

// func (m *MongoPkgStruct) ReConnect() error {
// 	// 古いコンテキストをキャンセル
// 	if m.Cancel != nil {
// 		m.Cancel()
// 	}

// 	// 新しい接続を確立
// 	client, ctx, cancelFunc, err := connect()
// 	if err != nil {
// 		return err
// 	}

// 	m.Ctx = ctx
// 	m.Cancel = cancelFunc
// 	m.Db = client.Database("chatapp")
// 	fmt.Println("ReConnected to MongoDB!")

// 	return nil
// }
