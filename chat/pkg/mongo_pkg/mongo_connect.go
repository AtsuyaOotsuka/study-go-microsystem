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
	Db     MongoDatabaseInterface
	Ctx    context.Context
	Cancel context.CancelFunc
}

type MongoPkgInterface interface {
	NewMongoConnect(database string) (*MongoPkgStruct, error)
}

type MongoPkg struct{}

func NewMongoPkg() *MongoPkg {
	return &MongoPkg{}
}

func (m *MongoPkg) connect() (*mongo.Client, context.Context, context.CancelFunc, error) {
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

func (m *MongoPkg) NewMongoConnect(database string) (*MongoPkgStruct, error) {
	client, ctx, cancelFunc, err := m.connect()
	if err != nil {
		return nil, err
	}

	mongoPkgStruct := &MongoPkgStruct{}
	mongoPkgStruct.Ctx = ctx
	mongoPkgStruct.Cancel = cancelFunc
	mongoClient := &RealMongoClient{client: client}
	mongoPkgStruct.Db = mongoClient.Database(database)
	fmt.Println("Connected to MongoDB!")

	return mongoPkgStruct, nil
}
