package mongo_dao

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ()

func DeleteById(cli *mongo.Collection, id any, opts ...*options.DeleteOptions) int64 {
	r, _ := cli.DeleteOne(context.TODO(), bson.M{"_id": id}, opts...)
	return r.DeletedCount
}
