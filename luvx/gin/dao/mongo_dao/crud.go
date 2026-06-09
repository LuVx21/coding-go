package mongo_dao

import (
	"context"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Distinct[T any](field string, filter any, opts ...options.Lister[options.DistinctOptions]) []T {
	if field == "" {
		return nil
	}
	distinctValues := WeiboFeedCol.Distinct(context.TODO(), field, filter, opts...)
	var r []T
	distinctValues.Decode(&r)
	return r
}

func IgnoreRetweet() []int64 {
	filter := bson.M{
		"groupId": 4670120389774996,
		"invalid": 0,
		"text":    "转发微博",
		"retweeted_status": bson.M{
			"$exists": true,
			"$ne":     nil,
		},
	}
	opts := options.Find().SetProjection(bson.M{"_id": 1, "retweeted_status": 1})
	ms, _ := mongodb.RowsMap(context.TODO(), WeiboFeedCol, filter, opts)

	ids, rids := sets.NewSet[int64](), sets.NewSet[int64]()
	for _, m := range *ms {
		if cell, ok := m["_id"]; ok {
			ids.Add(cast_x.ToInt64(cell))
		}
		if cell, ok := m["retweeted_status"]; ok {
			rids.Add(cast_x.ToInt64(cell))
		}
	}
	a, b := ids.ToSlice(), rids.ToSlice()
	if ids.Len() != 0 {
		WeiboFeedCol.UpdateMany(context.TODO(),
			bson.M{
				"_id":     bson.M{"$in": a},
				"invalid": 0,
			},
			bson.M{"$set": bson.M{"invalid": 1, "read": 1}},
		)
	}
	if rids.Len() != 0 {
		WeiboFeedCol.UpdateMany(context.TODO(),
			bson.M{
				"_id":     bson.M{"$in": b},
				"read":    0,
				"invalid": 1,
			},
			bson.M{"$set": bson.M{"invalid": 0}},
		)
	}
	return a
}
