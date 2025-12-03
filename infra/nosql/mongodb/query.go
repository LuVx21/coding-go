package mongodb

import (
	"context"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func First(cli *mongo.Collection, filter any, orderByKey string) (bson.M, error) {
	return firstOrLast(cli, filter, orderByKey, true)
}
func Last(cli *mongo.Collection, filter any, orderByKey string) (bson.M, error) {
	return firstOrLast(cli, filter, orderByKey, false)
}
func firstOrLast(cli *mongo.Collection, filter any, orderByKey string, asc bool) (bson.M, error) {
	opts := options.FindOne().SetSort(bson.M{orderByKey: common_x.IfThen(asc, 1, -1)})
	var m bson.M
	e := cli.FindOne(context.TODO(), filter, opts).Decode(&m)
	return m, e
}
