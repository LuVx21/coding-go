package mongodb

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertMany(ctx context.Context, col *mongo.Collection, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if many, err := col.InsertMany(ctx, documents, opts...); err == nil {
		return many, err
	}
	for _, document := range documents {
		_, _ = col.InsertOne(ctx, document)
	}
	return nil, nil
}

func RowsMap(ctx context.Context, col *mongo.Collection, filter any, opts ...*options.FindOptions) (*[]bson.M, error) {
	cur, err := col.Find(ctx, filter, opts...)
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error closing cursor")
		}
	}(cur, ctx)

	if err != nil {
		return nil, err
	}
	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
