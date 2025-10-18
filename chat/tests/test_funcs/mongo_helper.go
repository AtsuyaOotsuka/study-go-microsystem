package test_funcs

import (
	"context"
	"fmt"
	"microservices/chat/internal/model"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetUpMongoTestDatabase() *TestMongoStruct {
	ctx := context.Background()
	uri := os.Getenv("MONGODB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := client.Database("chatapp")
	return &TestMongoStruct{
		DB:     db,
		Ctx:    ctx,
		Client: client,
	}
}

type TestMongoStruct struct {
	DB     *mongo.Database
	Ctx    context.Context
	Client *mongo.Client
}

func (m *TestMongoStruct) Disconnect() error {
	fmt.Println("Disconnecting MongoDB client...")
	return m.Client.Disconnect(m.Ctx)
}

func (m *TestMongoStruct) MongoCleanUp() error {

	var err error

	err = m.DB.Collection(model.RoomCollectionName).Drop(m.Ctx)
	if err != nil {
		return err
	}
	err = m.DB.Collection(model.ChatMessageCollectionName).Drop(m.Ctx)
	if err != nil {
		return err
	}

	fmt.Println("MongoDB cleaned up for tests.")
	return nil
}

func (m *TestMongoStruct) ExistContents(collectionName string, filter interface{}) (bool, error) {
	count, err := m.DB.Collection(collectionName).CountDocuments(m.Ctx, filter)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (m *TestMongoStruct) CountContents(collectionName string, filter interface{}) (int64, error) {
	count, err := m.DB.Collection(collectionName).CountDocuments(m.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
