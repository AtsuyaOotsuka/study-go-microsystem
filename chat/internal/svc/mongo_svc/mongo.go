package mongo_svc

import (
	"microservices/chat/internal/model"
	"microservices/chat/pkg/mongo_pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoSvcInterface interface {
	CreateRoom(room model.Room, mongo_pkg mongo_pkg.MongoPkgInterface) (interface{}, error)
	GetRoomByID(roomID string, mongo_pkg mongo_pkg.MongoPkgInterface) (model.Room, error)
	JoinRoom(roomID string, userID int, mongo_pkg mongo_pkg.MongoPkgInterface) error
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

func (m *MongoSvcStruct) GetRoomByID(roomID string, mongo_pkg mongo_pkg.MongoPkgInterface) (model.Room, error) {
	mongo, err := Init(mongo_pkg)
	if err != nil {
		return model.Room{}, err
	}

	defer mongo.MongoPkgStruct.Cancel()

	collection := mongo.MongoPkgStruct.Db.Collection(model.RoomCollectionName)

	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return model.Room{}, err
	}

	var room model.Room
	err = collection.FindOne(mongo.MongoPkgStruct.Ctx, bson.M{"_id": id}, &room)
	if err != nil {
		return model.Room{}, err
	}

	return room, nil
}

func (m *MongoSvcStruct) JoinRoom(roomID string, userID int, mongo_pkg mongo_pkg.MongoPkgInterface) error {
	mongo, err := Init(mongo_pkg)
	if err != nil {
		return err
	}

	defer mongo.MongoPkgStruct.Cancel()

	collection := mongo.MongoPkgStruct.Db.Collection(model.RoomCollectionName)

	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(
		mongo.MongoPkgStruct.Ctx,
		bson.M{"_id": id},
		bson.M{"$addToSet": bson.M{"members": userID}}, // 重複追加防止してくれる
	)
	if err != nil {
		return err
	}

	return nil
}
