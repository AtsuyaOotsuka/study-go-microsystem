package mongo_pkg

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo Cursor
type MongoCursorInterface interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Close(ctx context.Context) error
}

type RealMongoCursor struct {
	cursor *mongo.Cursor
}

func (r *RealMongoCursor) Next(ctx context.Context) bool {
	return r.cursor.Next(ctx)
}

func (r *RealMongoCursor) Decode(val interface{}) error {
	return r.cursor.Decode(val)
}

func (r *RealMongoCursor) Close(ctx context.Context) error {
	return r.cursor.Close(ctx)
}

// Mongo Database
type MongoDatabaseInterface interface {
	Collection(name string) MongoCollectionInterface
}

type RealMongoDatabase struct {
	db *mongo.Database
}

func (r *RealMongoDatabase) Collection(name string) MongoCollectionInterface {
	// *mongo.Collection を RealMongoCollection に包んで返す
	return &RealMongoCollection{coll: r.db.Collection(name)}
}

// Mongo Collection
type MongoCollectionInterface interface {
	InsertOne(ctx context.Context, document interface{}) (string, error)
	Find(ctx context.Context, filter interface{}) (cursor MongoCursorInterface, err error)
	FindOne(ctx context.Context, filter interface{}, object interface{}) error
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

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

func (r *RealMongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.coll.UpdateMany(ctx, filter, update)
}

func (r *RealMongoCollection) Find(ctx context.Context, filter interface{}) (cursor MongoCursorInterface, err error) {
	cursor, err = r.coll.Find(ctx, filter)
	return
}

func (r *RealMongoCollection) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	err := r.coll.FindOne(ctx, filter).Decode(object)
	return err
}

func (r *RealMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.coll.UpdateOne(ctx, filter, update)
}

func (r *RealMongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return r.coll.DeleteOne(ctx, filter)
}

// Mongo Client
type RealMongoClient struct {
	client *mongo.Client
}

func (r *RealMongoClient) Database(name string) MongoDatabaseInterface {
	return &RealMongoDatabase{db: r.client.Database(name)}
}
