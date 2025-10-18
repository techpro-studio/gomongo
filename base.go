package gomongo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (m *BaseMongoRepository) UpdateOne(ctx context.Context, q bson.M, update bson.M) error {
	result, err := m.Collection().UpdateOne(ctx, q, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount != 1 {
		return errors.New(fmt.Sprintf("Incorrect modified count should be 1 got %d", result.ModifiedCount))
	}
	return err
}

func (m *BaseMongoRepository) DeleteOne(ctx context.Context, q bson.M) error {
	result, err := m.Collection().DeleteOne(ctx, q)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Incorrect deleted count should be 1 got %d", result.DeletedCount))
	}
	return err
}

func (m *BaseMongoRepository) InsertOne(ctx context.Context, newValue interface{}) (*bson.ObjectID, error) {
	result, err := m.Collection().InsertOne(ctx, newValue)
	if err != nil {
		return nil, err
	}
	objID := result.InsertedID.(bson.ObjectID)
	return &objID, nil
}

func (m *BaseMongoRepository) GetList(ctx context.Context, result interface{}, q bson.M, skip, limit *int, sort *bson.D) (int, error) {

	opts := options.Find()
	if skip != nil {
		opts.SetSkip(int64(*skip))
	}
	if limit != nil {
		opts.SetSkip(int64(*limit))
	}
	if sort != nil {
		opts.SetSort(*sort)
	}

	cursor, err := m.Collection().Find(ctx, q, opts)
	if err != nil {
		return -1, err
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return -1, err
	}
	count, err := m.Collection().CountDocuments(ctx, q)
	if err != nil {
		return -1, err
	}
	return int(count), nil
}
