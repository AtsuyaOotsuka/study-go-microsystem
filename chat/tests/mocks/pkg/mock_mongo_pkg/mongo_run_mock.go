package mock_mongo_pkg

import (
	"context"
	"microservices/chat/pkg/mongo_pkg"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCursorMock struct {
	mock.Mock
}

func (m *MongoCursorMock) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MongoCursorMock) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

func (m *MongoCursorMock) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MongoDatabaseMock struct {
	mock.Mock
}

func (m *MongoDatabaseMock) Collection(name string) mongo_pkg.MongoCollectionInterface {
	args := m.Called(name)
	return args.Get(0).(mongo_pkg.MongoCollectionInterface)
}

type MongoCollectionMock struct {
	mock.Mock
}

func (m *MongoCollectionMock) InsertOne(ctx context.Context, document interface{}) (string, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(string), args.Error(1)
}

func (m *MongoCollectionMock) Find(ctx context.Context, filter interface{}) (cursor mongo_pkg.MongoCursorInterface, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(mongo_pkg.MongoCursorInterface), args.Error(1)
}

func (m *MongoCollectionMock) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	args := m.Called(ctx, filter, object)
	return args.Error(0)
}

func (m *MongoCollectionMock) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MongoCollectionMock) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

type MongoPkgStructMock struct {
	Ctx context.Context
	Db  mongo_pkg.MongoDatabaseInterface
}

type MongoMock struct {
	MongoPkgStruct *MongoPkgStructMock
}

type MongoCollectionInsertErrorMock struct {
	mock.Mock
}

func (m *MongoCollectionInsertErrorMock) InsertOne(ctx context.Context, document interface{}) (string, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(string), args.Error(1)
}

func (m *MongoCollectionInsertErrorMock) Find(ctx context.Context, filter interface{}) (cursor mongo_pkg.MongoCursorInterface, err error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(mongo_pkg.MongoCursorInterface), args.Error(1)
}

func (m *MongoCollectionInsertErrorMock) FindOne(ctx context.Context, filter interface{}, object interface{}) error {
	args := m.Called(ctx, filter, object)
	return args.Error(0)
}

func (m *MongoCollectionInsertErrorMock) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MongoCollectionInsertErrorMock) UpdateMany(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}
