package mongodb

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InsertMany(ctx context.Context, col *mongo.Collection, documents []any, opts ...options.Lister[options.InsertManyOptions]) (*mongo.InsertManyResult, error) {
	if many, err := col.InsertMany(ctx, documents, opts...); err == nil {
		return many, err
	}
	for _, document := range documents {
		_, _ = col.InsertOne(ctx, document)
	}
	return nil, nil
}

func RowsMap(ctx context.Context, col *mongo.Collection, filter any, opts ...options.Lister[options.FindOptions]) (*[]bson.M, error) {
	cur, err := col.Find(ctx, filter, opts...)
	if err != nil {
		slog.Warn("mongo查询错误", "err", err.Error())
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error closing cursor")
		}
	}(cur, ctx)

	var results []bson.M
	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func MD(m bson.M) bson.D {
	r := make(bson.D, 0, len(m))
	for k, v := range m {
		r = append(r, bson.E{Key: k, Value: v})
	}
	return r
}

func DM(es bson.D) bson.M {
	r := make(bson.M, len(es))
	for i := range es {
		e := es[i]
		r[e.Key] = e.Value
	}
	return r
}
func FindInD(es bson.D, name string) any {
	for i := range es {
		e := es[i]
		if e.Key == name {
			return e.Value
		}
	}
	return nil
}
