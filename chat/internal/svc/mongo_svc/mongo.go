package mongo_svc

import (
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"
)

type MongoSvcInterface interface {
	CreateRoom(room model.Room, mongo_pkg mongo_pkg.MongoPkgInterface) (interface{}, error)
}

type MongoSvcStruct struct {
	Db mongo_pkg.MongoDatabaseInterface
}

func NewMongoSvc(db mongo_pkg.MongoDatabaseInterface) *MongoSvcStruct {
	return &MongoSvcStruct{
		Db: db,
	}
}

type Mongo struct {
	MongoPkgStruct *mongo_pkg.MongoPkgStruct
}

func Init(mongo_pkg mongo_pkg.MongoPkgInterface) (mongo *Mongo, err error) {
	mongoPkgStruct, err := mongo_pkg.NewMongoConnect("chatapp")
	if err != nil {
		return nil, err
	}
	mongo = &Mongo{
		MongoPkgStruct: mongoPkgStruct,
	}
	return mongo, nil
}

func (m *MongoSvcStruct) CreateRoom(room model.Room, mongo_pkg mongo_pkg.MongoPkgInterface) (interface{}, error) {
	mongo, err := Init(mongo_pkg)
	if err != nil {
		return nil, err
	}

	defer mongo.MongoPkgStruct.Cancel()

	collection := mongo.MongoPkgStruct.Db.Collection(model.RoomCollectionName)

	InsertedID, err := collection.InsertOne(mongo.MongoPkgStruct.Ctx, room)
	if err != nil {
		return nil, err
	}

	return InsertedID, nil
}
