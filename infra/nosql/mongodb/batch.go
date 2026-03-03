package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InsertOneBatch(clients []*mongo.Collection, ctx context.Context, document any, opts ...options.Lister[options.InsertOneOptions]) {
	for _, client := range clients {
		go func() {
			_, _ = client.InsertOne(ctx, document, opts...)
		}()
	}
}

func InsertManyBatch(clients []*mongo.Collection, ctx context.Context, documents []any, opts ...options.Lister[options.InsertManyOptions]) {
	for _, client := range clients {
		go func() {
			_, _ = client.InsertMany(ctx, documents, opts...)
		}()
	}
}
