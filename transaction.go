package gomongo

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// In case you need to return no Value

type Void struct{}

type TransactionFunction[T any] func(sc context.Context) (T, error)

func InTransactionSession[T any](ctx context.Context, client *mongo.Client, f TransactionFunction[T]) (T, error) {
	session, err := client.StartSession()
	if err != nil {
		var zero T
		return zero, err
	}
	defer session.EndSession(ctx)
	var value T
	err = mongo.WithSession(ctx, session, func(sc context.Context) error {
		err := session.StartTransaction()
		if err != nil {
			return err
		}
		value, err = f(sc)
		if err != nil {
			return err
		}
		return session.CommitTransaction(sc)
	})
	return value, err
}
