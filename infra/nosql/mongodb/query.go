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

func Exist(cli *mongo.Collection, filter any) (bool, error) {
	c, err := cli.CountDocuments(context.TODO(), filter)
	return c > 0, err
}
func Exists(cli *mongo.Collection, ids ...any) (map[any]struct{}, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	ms, err := RowsMap(context.TODO(), cli, filter, opts)
	if err != nil {
		return nil, err
	}
	a := make(map[any]struct{}, len(*ms))
	for i := range *ms {
		m := (*ms)[i]
		a[m["_id"]] = struct{}{}
	}
	return a, nil
}
