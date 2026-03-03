package mongo_dao

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ()

func DeleteById(cli *mongo.Collection, id any, opts ...options.Lister[options.DeleteOneOptions]) int64 {
	r, _ := cli.DeleteOne(context.TODO(), bson.M{"_id": id}, opts...)
	return r.DeletedCount
}
