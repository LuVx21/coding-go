package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/iterators"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_json = `
{
	"source": {
		"uri": "",
		"database": "boot",
		"collection":"user1"
	},
	"sink": {
		"uri": "",
		"database": "boot",
		"collection":"user2"
	},
	"filter": [
		{
			"key": "invalid",
			"operator": "eq",
			"value": 0
		}
	],
	"uniqueField": "_id",
	"order": "asc",
	"batchSize": 500,
	"onSuccess": "markDel",
	"compareField": []
}
`
)

var (
	onSuccessValid = []string{"none", "delete"}
)

type config struct {
	Source, Sink       mongoConfig
	Filter             []kv
	BatchSize          int64
	OnSuccess          string
	UniqueField, Order string // 仅支持数字类型唯一性字段, asc时顺序,其他逆序
	CompareField       []string
}
type mongoConfig struct {
	Uri, Database, Collection string
}
type kv struct {
	Key      string
	Operator string
	Value    any
}

func main() {
	var _config config
	err := json.Unmarshal([]byte(_json), &_config)
	if err != nil {
		return
	}
	source, sink, uk, desc := _config.Source, _config.Sink, _config.UniqueField, _config.Order != "asc"
	sort := bson.M{uk: common_x.IfThen(desc, -1, 1)}

	sourceClient, sourceErr := mongo.Connect(context.TODO(), options.Client().ApplyURI(source.Uri))
	if sourceErr != nil {
		slog.Error("mongo source 连接失败")
		return
	}
	sinkClient, sinkErr := mongo.Connect(context.TODO(), options.Client().ApplyURI(sink.Uri))
	if sinkErr != nil {
		slog.Error("mongo sink 连接失败")
		return
	}
	sourceDb, sinkDB := sourceClient.Database(source.Database).Collection(source.Collection), sinkClient.Database(sink.Database).Collection(sink.Collection)

	filter := buildFilter(_config.Filter)
	opts := options.FindOne().SetSort(sort)
	var minMax bson.M
	_ = sourceDb.FindOne(context.TODO(), filter, opts).Decode(&minMax)
	cursor := cast_x.ToInt64(minMax[uk]) + common_x.IfThen(desc, int64(1), int64(-1))
	iterator := iterators.NewCursorIteratorSimple(
		cursor, false,
		func(_cursor int64) []bson.M {
			_filter := append(filter, bson.E{Key: uk, Value: bson.M{common_x.IfThen(desc, "$lt", "$gt"): cursor}})
			opts := options.Find().SetSort(sort).SetLimit(_config.BatchSize)
			rowsMap, _ := mongodb.RowsMap(context.Background(), sourceDb, _filter, opts)
			return *rowsMap
		},
		func(curId int64, items []bson.M) int64 {
			fmt.Println("---------------------------------------------------------------")
			if int64(len(items)) < _config.BatchSize {
				return -1
			}
			cursor = cast_x.ToInt64(items[len(items)-1][uk])
			return cursor
		},
		func(i int64) bool {
			return i <= 0
		},
	)
	iterator.ForEachRemaining(func(m bson.M) {
		_, e := sinkDB.InsertOne(context.TODO(), m)
		if e == nil {
			switch _config.OnSuccess {
			case "delete":
				sourceDb.DeleteOne(context.TODO(), bson.M{uk: m[uk]})
			case "markDel":
				sourceDb.UpdateOne(context.TODO(), bson.M{uk: m[uk]}, bson.D{{Key: "$set", Value: bson.D{{Key: "markDel", Value: 1}}}})
			}
			return
		}

		if strings.Contains(e.Error(), "duplicate key error collection") {
			if len(_config.CompareField) > 0 {
				var target bson.M
				sinkDB.FindOne(context.Background(), bson.M{uk: m[uk]}).Decode(&target)
				update := make(bson.D, 0)
				for _, f := range _config.CompareField {
					origina := m[f]
					if origina != target[f] {
						update = append(update, bson.E{Key: f, Value: origina})
					}
				}
				if len(update) > 0 {
					sinkDB.UpdateOne(context.Background(), bson.M{uk: m[uk]}, bson.D{{Key: "$set", Value: update}})
				}
			}
			return
		}

		slog.Error("同步写入失败", "错误", e.Error(), "id", m[uk])
	})
}

func buildFilter(filters []kv) bson.D {
	filter := make(bson.D, 0)

	for _, kv := range filters {
		v := kv.Value
		switch kv.Operator {
		case "eq":
		case "ne":
			v = bson.M{"$ne": v}
		case "exists":
			v = bson.M{"$exists": v}
		case "in":
			v = bson.M{"$in": v}
		case "nin":
			v = bson.M{"$nin": v}
		case "regex":
		case "size":
		default:
		}
		filter = append(filter, bson.E{Key: kv.Key, Value: v})
	}

	return filter
}
