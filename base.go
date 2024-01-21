package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseMongoRepository struct {
	database       *mongo.Database
	collectionName string
}

func NewBaseMongoRepository(database *mongo.Database, collectionName string) *BaseMongoRepository {
	return &BaseMongoRepository{database: database, collectionName: collectionName}
}

func (m *BaseMongoRepository) Collection() *mongo.Collection {
	return m.database.Collection(m.collectionName)
}

func (m *BaseMongoRepository) UpdateOne(ctx context.Context, q bson.M, update bson.M) bool {
	result, err := m.Collection().UpdateOne(ctx, q, update)
	if err != nil {
		panic(err)
	}
	return result.ModifiedCount == 1
}

func (m *BaseMongoRepository) DeleteOne(ctx context.Context, q bson.M) bool {
	result, err := m.Collection().DeleteOne(ctx, q)
	if err != nil {
		panic(err)
	}
	return result.DeletedCount == 1
}

func (m *BaseMongoRepository) InsertOne(ctx context.Context, newValue interface{}) primitive.ObjectID {
	result, err := m.Collection().InsertOne(ctx, newValue)
	if err != nil {
		panic(err)
	}
	return result.InsertedID.(primitive.ObjectID)
}

func (m *BaseMongoRepository) GetList(ctx context.Context, result interface{}, q bson.M, skip, limit *int, sort *bson.D) int {
	cursor, err := m.Collection().Find(ctx, q, &options.FindOptions{
		Skip:  Int64Ptr(skip),
		Limit: Int64Ptr(limit),
		Sort:  sort,
	})
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, result)
	if err != nil {
		panic(err)
	}
	count, err := m.Collection().CountDocuments(ctx, q)
	if err != nil {
		panic(err)
	}
	return int(count)
}

func InTransactionSession(Client *mongo.Client, ctx context.Context, f func(sc mongo.SessionContext) error) {
	var session mongo.Session
	session, err := Client.StartSession()
	if err != nil {
		panic(err)
	}
	err = session.StartTransaction()
	if err != nil {
		panic(err)
	}
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		err := f(sc)
		if err != nil {
			return err
		}
		err = session.CommitTransaction(sc)
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	session.EndSession(ctx)
}
