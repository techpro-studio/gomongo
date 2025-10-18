package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

func StrToObjId(id *string) *bson.ObjectID {
	if id == nil {
		return nil
	}
	objID, err := bson.ObjectIDFromHex(*id)
	if err != nil {
		panic(err)
	}
	return &objID
}

func ConnectToMongo(uri string) *mongo.Client {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	return client
}

func ObjIdToStr(id *bson.ObjectID) *string {
	if id == nil {
		return nil
	}
	hex := (*id).Hex()
	return &hex
}

func ObjIdListToStrList(list []bson.ObjectID) []string {
	var result []string
	for _, item := range list {
		result = append(result, item.Hex())
	}
	return result
}

func StrListToObjIdList(list []string) []bson.ObjectID {
	var result []bson.ObjectID
	for _, item := range list {
		result = append(result, *StrToObjId(&item))
	}
	return result
}
