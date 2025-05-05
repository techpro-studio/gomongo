package gomongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TypedRepository[T any, M ModelConverted[T]] struct {
	BaseMongoRepository
}

func NewTypedRepository[T any, M ModelConverted[T]](database *mongo.Database, collectionName string) *TypedRepository[T, M] {
	return &TypedRepository[T, M]{BaseMongoRepository: *NewBaseMongoRepository(database, collectionName)}
}

func (r *TypedRepository[T, M]) GetOne(ctx context.Context, q bson.M) (*T, error) {
	var schema M
	err := r.Collection().FindOne(ctx, q).Decode(&schema)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return schema.ToModel(), nil
}

func (r *TypedRepository[T, M]) GetByIdList(ctx context.Context, idList []string) ([]T, error) {
	return r.GetTypedList(ctx, bson.M{"_id": bson.M{"$in": StrListToObjIdList(idList)}}, nil, nil, nil)
}

func (r *TypedRepository[T, M]) GetTypedList(ctx context.Context, q bson.M, skip, limit *int, sort *bson.D) ([]T, error) {
	var schema []M
	_, err := r.GetList(ctx, &schema, q, skip, limit, sort)
	if err != nil {
		return nil, err
	}
	return SliceMap(schema, func(input M) T {
		return *input.ToModel()
	}), err
}

func (r *TypedRepository[T, M]) GetOneById(ctx context.Context, id string) (*T, error) {
	return r.GetOne(ctx, bson.M{"_id": *StrToObjId(&id)})
}

func (r *TypedRepository[T, M]) GetPaginatedList(ctx context.Context, q bson.M, after *string, limit *int) (*PaginatedList[T], error) {
	query := NewFromListQuery(limit, after).SetDefaultSortById()
	query.Query = q
	return GetPaginatedListForQuery[T, M](ctx, r.Collection(), *query)
}
