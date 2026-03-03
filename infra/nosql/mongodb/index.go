package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetAllIndex(cli *mongo.Collection) map[string]bson.M {
	ctx := context.Background()

	cursor, err := cli.Indexes().List(ctx)
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	r := make(map[string]bson.M, 0)
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil
		}
		r[result["name"].(string)] = result
	}

	return r
}
