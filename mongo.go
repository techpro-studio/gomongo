package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func StrToObjId(id *string) *primitive.ObjectID {
	if id == nil {
		return nil
	}
	objID, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		panic(err)
	}
	return &objID
}

func ConnectToMongo(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	return client
}

func ObjIdToStr(id *primitive.ObjectID) *string {
	if id == nil {
		return nil
	}
	hex := (*id).Hex()
	return &hex
}

func ObjIdListToStrList(list []primitive.ObjectID) []string {
	var result []string
	for _, item := range list {
		result = append(result, item.Hex())
	}
	return result
}

func StrListToObjIdList(list []string) []primitive.ObjectID {
	var result []primitive.ObjectID
	for _, item := range list {
		result = append(result, *StrToObjId(&item))
	}
	return result
}
