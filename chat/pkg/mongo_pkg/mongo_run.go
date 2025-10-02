package mongo_pkg

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDatabaseInterface interface {
	Collection(name string) MongoCollectionInterface
}

type MongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}) (string, error)
	FindOne(ctx context.Context, filter interface{}, object interface{}) error
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
}

// RealMongoDatabase wraps *mongo.Database
type RealMongoDatabase struct {
	db *mongo.Database
}

// 正しく interface を返すようにする
func (r *RealMongoDatabase) Collection(name string) MongoCollectionInterface {
	// *mongo.Collection を RealMongoCollection に包んで返す
	return &RealMongoCollection{coll: r.db.Collection(name)}
}

// RealMongoCollection wraps *mongo.Collection
type RealMongoCollection struct {
	coll *mongo.Collection
}

func (r *RealMongoCollection) InsertOne(ctx context.Context, document interface{}) (string, error) {
	result, err := r.coll.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertedID is not an ObjectID: %#v", result.InsertedID)
	}

	return id.Hex(), nil
}

func (r *RealMongoCollection) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	err := r.coll.FindOne(ctx, filter).Decode(object)
	return err
}

func (r *RealMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.coll.UpdateOne(ctx, filter, update)
}

type RealMongoClient struct {
	client *mongo.Client
}

func (r *RealMongoClient) Database(name string) MongoDatabaseInterface {
	return &RealMongoDatabase{db: r.client.Database(name)}
}
