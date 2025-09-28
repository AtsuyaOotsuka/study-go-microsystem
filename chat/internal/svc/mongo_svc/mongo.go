package mongo_svc

import (
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"
)

type MongoSvcInterface interface {
	CreateRoom(room model.Room) (interface{}, error)
}

type MongoSvcStruct struct{}

func NewMongoSvc() *MongoSvcStruct {
	return &MongoSvcStruct{}
}

type Mongo struct {
	MongoPkgStruct *mongo_pkg.MongoPkgStruct
}

func Init() (mongo *Mongo, err error) {
	mongoPkgStruct, err := mongo_pkg.NewMongoConnect("chatapp")
	if err != nil {
		return nil, err
	}
	mongo = &Mongo{
		MongoPkgStruct: mongoPkgStruct,
	}
	return mongo, nil
}

func (m *MongoSvcStruct) CreateRoom(room model.Room) (interface{}, error) {
	mongo, err := Init()
	if err != nil {
		return nil, err
	}

	defer mongo.MongoPkgStruct.Cancel()

	collection := mongo.MongoPkgStruct.Db.Collection(model.RoomCollectionName)

	result, err := collection.InsertOne(mongo.MongoPkgStruct.Ctx, room)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
