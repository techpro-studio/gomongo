package gomongo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

const (
	KEnvMongoURL           = "MONGO_URL"
	KEnvMongoMaxPoolSize   = "MONGO_MAX_POOL_SIZE"
	KEnvMongoMinPoolSize   = "MONGO_MIN_POOL_SIZE"
	KEnvMongoMaxConnecting = "MONGO_MAX_CONNECTING"
)

type MongoConnectionParameters struct {
	Uri           string
	MaxPoolSize   uint64
	MinPoolSize   uint64
	MaxConnecting uint64
}

func RequiredError(key string) error {
	return fmt.Errorf("%s is required", key)
}

func ParseMongoConnectionParameterFromEnv(prefix string) (*MongoConnectionParameters, error) {
	buildKey := func(key string) string {
		if prefix == "" {
			return key
		}
		return fmt.Sprintf("%s_%s", prefix, key)
	}

	mainMongoURL := os.Getenv(buildKey(KEnvMongoURL))
	if mainMongoURL == "" {
		return nil, RequiredError(buildKey(KEnvMongoURL))
	}

	maxPoolSize, err := strconv.ParseUint(os.Getenv(buildKey(KEnvMongoMaxPoolSize)), 10, 64)
	if err != nil {
		maxPoolSize = 100
	}

	minPoolSize, err := strconv.ParseUint(os.Getenv(buildKey(KEnvMongoMinPoolSize)), 10, 64)
	if err != nil {
		minPoolSize = 0
	}

	maxConnecting, err := strconv.ParseUint(os.Getenv(buildKey(KEnvMongoMaxConnecting)), 10, 64)
	if err != nil {
		maxConnecting = 2
	}

	return &MongoConnectionParameters{
		Uri:           mainMongoURL,
		MaxPoolSize:   maxPoolSize,
		MinPoolSize:   minPoolSize,
		MaxConnecting: maxConnecting,
	}, nil
}

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

func ConnectToMongo(params *MongoConnectionParameters) *mongo.Client {
	opts := options.Client().
		ApplyURI(params.Uri).
		SetBSONOptions(&options.BSONOptions{ObjectIDAsHexString: true}).
		SetMaxPoolSize(params.MaxPoolSize).
		SetMinPoolSize(params.MinPoolSize).
		SetMaxConnecting(params.MaxConnecting)

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
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
