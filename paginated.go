package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQuery struct {
	Limit       *int64
	Query       bson.M
	AfterQuery  bson.M
	Sort        *bson.D
	IgnoreTotal bool
}

func NewMongoQuery() *MongoQuery {
	return &MongoQuery{
		AfterQuery: bson.M{},
		Query:      bson.M{},
	}
}

func (q *MongoQuery) SetDefaultSortById() *MongoQuery {
	sort := bson.D{{Key: "_id", Value: -1}}
	q.Sort = &sort
	return q
}

func (q *MongoQuery) SetSort(sort bson.D) *MongoQuery {
	q.Sort = &sort
	return q
}

func (q *MongoQuery) SetLimit(limit int) *MongoQuery {
	q.Limit = Int64Ptr(&limit)
	return q
}

func (q *MongoQuery) SetIgnoreTotal(flag bool) *MongoQuery {
	q.IgnoreTotal = flag
	return q
}

func (q *MongoQuery) SetQuery(query bson.M) *MongoQuery {
	q.Query = query
	return q
}

func (q *MongoQuery) SetAfterQuery(afterQuery bson.M) *MongoQuery {
	q.AfterQuery = afterQuery
	return q
}

func NewFromListQuery(limit *int, afterId *string) *MongoQuery {
	afterQuery := bson.M{}

	if afterId != nil {
		afterQuery["_id"] = bson.M{"$lt": *StrToObjId(afterId)}
	}

	return &MongoQuery{
		Limit:      Int64Ptr(limit),
		Query:      bson.M{},
		AfterQuery: afterQuery,
	}
}

type PaginatedList[T any] struct {
	HasMore bool `json:"has_more"`
	Total   int  `json:"total"`
	Items   []T  `json:"items"`
}

type ModelConverted[M any] interface {
	ToModel() *M
}

func GetPaginatedListForQuery[Model any, Schema ModelConverted[Model]](ctx context.Context, collection *mongo.Collection, mongoQuery MongoQuery) (*PaginatedList[Model], error) {
	var schema []Schema
	var limitPtr *int
	if mongoQuery.Limit != nil {
		var limit = int(*mongoQuery.Limit) + 1
		limitPtr = &limit
	}

	opts := options.FindOptions{
		Limit: Int64Ptr(limitPtr),
	}
	if mongoQuery.Sort != nil {
		opts.Sort = mongoQuery.Sort
	}
	fullQuery := bson.M{}

	for k, v := range mongoQuery.Query {
		fullQuery[k] = v
	}
	for k, v := range mongoQuery.AfterQuery {
		fullQuery[k] = v
	}

	cursor, err := collection.Find(ctx, fullQuery, &opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &schema)
	if err != nil {
		return nil, err
	}

	var total int64
	if !mongoQuery.IgnoreTotal {
		total, err = collection.CountDocuments(ctx, mongoQuery.Query)
		if err != nil {
			return nil, err
		}
	}

	hasMore := limitPtr != nil && len(schema) == *limitPtr

	if hasMore {
		schema = schema[:len(schema)-1]
	}

	return &PaginatedList[Model]{
		HasMore: hasMore,
		Total:   int(total),
		Items: SliceMap(schema, func(input Schema) Model {
			return *input.ToModel()
		}),
	}, nil
}
