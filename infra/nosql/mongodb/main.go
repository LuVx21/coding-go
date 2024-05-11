package mongodb

import (
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log/slog"
)

func RowsMap(ctx context.Context, col *mongo.Collection, filter interface{}, opts ...*options.FindOptions) (*[]bson.M, error) {
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
