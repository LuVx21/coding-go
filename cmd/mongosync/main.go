package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/iterators"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	"batchSize": 1024,
	"onSuccess": "",
	"compareField": [],
	"batchMove": true,
	"batchMoveSize": 1024
}
`
)

var (
	onSuccessValid = []string{"none", "delete"} //lint:ignore U1000 忽略
)

type config struct {
	Source, Sink       mongoConfig
	Filter             []kv
	BatchSize          int64
	OnSuccess          string
	UniqueField, Order string // 仅支持数字类型唯一性字段, asc时顺序,其他逆序
	CompareField       []string
	BatchMove          bool
	BatchMoveSize      int
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
		slog.Error("加载配置异常", "err", err.Error())
		return
	}
	source, sink, uk, desc := _config.Source, _config.Sink, _config.UniqueField, _config.Order != "asc"
	sort := bson.M{uk: common_x.IfThen(desc, -1, 1)}

	sourceClient, sourceErr := mongo.Connect(options.Client().ApplyURI(source.Uri))
	if sourceErr != nil {
		slog.Error("mongo source 连接失败", "错误", sourceErr.Error())
		return
	}
	sinkClient, sinkErr := mongo.Connect(options.Client().ApplyURI(sink.Uri))
	if sinkErr != nil {
		slog.Error("mongo sink 连接失败", "错误", sinkErr.Error())
		return
	}
	sourceDb := sourceClient.Database(source.Database).Collection(source.Collection)
	sinkDB := sinkClient.Database(sink.Database).Collection(common_x.IfThen(len(sink.Collection) == 0, source.Collection, sink.Collection))

	existIndex := mongodb.GetAllIndex(sinkDB)
	indexs, _ := sourceDb.Indexes().ListSpecifications(context.TODO())
	for _, index := range indexs {
		if index.Name == "_id_" {
			continue
		}
		_, exist := existIndex[index.Name]
		if exist {
			continue
		}
		sinkDB.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
			Keys: index.KeysDocument,
			Options: options.Index().
				SetName(index.Name).
				SetUnique(index.Unique != nil && *index.Unique),
		})
	}

	filter := buildFilter(_config.Filter)
	opts := options.FindOne().SetSort(sort)
	var minMax bson.M
	_ = sourceDb.FindOne(context.TODO(), filter, opts).Decode(&minMax)
	cursor := cast_x.ToInt64(minMax[uk]) + common_x.IfThen(desc, int64(1), int64(-1))
	iterator := iterators.NewCursorIteratorSimple(
		cursor, false,
		func(_cursor int64) []bson.M {
			_filter := append(filter, bson.E{Key: uk, Value: bson.M{common_x.IfThen(desc, "$lt", "$gt"): _cursor}})
			opts := options.Find().SetSort(sort).SetLimit(_config.BatchSize)
			rowsMap, _ := mongodb.RowsMap(context.Background(), sourceDb, _filter, opts)
			return *rowsMap
		},
		func(curId int64, items []bson.M) int64 {
			fmt.Println("---------------------------------------------------------------")
			if int64(len(items)) < _config.BatchSize {
				return -1
			}
			return cast_x.ToInt64(items[len(items)-1][uk])
		},
		func(i int64) bool {
			return i <= 0
		},
	)

	if _config.BatchMove {
		batch(iterator, sourceDb, sinkDB, _config)
	} else {
		iterator.ForEachRemaining(func(m bson.M) {
			_, e := sinkDB.InsertOne(context.TODO(), m)
			if e == nil {
				postMoveSuccess([]any{m[uk]}, sourceDb, _config)
				return
			}
			postMoveError(e, []any{m}, sinkDB, _config)
		})
	}
}

func batch(iterator *iterators.CursorIterator[bson.M, int64, []bson.M], sourceDb, sinkDB *mongo.Collection, _config config) {
	temp := make([]any, 0)

	iterator.ForEachRemaining(func(m bson.M) {
		temp = append(temp, m)
		if !iterator.HasNext() || len(temp) >= _config.BatchMoveSize {
			imr, e := sinkDB.InsertMany(context.TODO(), temp, options.InsertMany().SetOrdered(false))
			if e != nil {
				postMoveError(e, temp, sinkDB, _config)
				return
			}
			if len(imr.InsertedIDs) > 0 {
				slog.Info("批量迁移", "移动数量:", len(imr.InsertedIDs))
				postMoveSuccess(imr.InsertedIDs, sourceDb, _config)
			}
			temp = temp[:0]
		}
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

func postMoveError(e error, datas []any, sinkDB *mongo.Collection, _config config) {
	if e == nil {
		return
	}

	wes := make([]mongo.WriteError, 0)
	// insertOne时出现
	if a, ok := e.(mongo.WriteException); ok {
		wes = append(wes, a.WriteErrors...)
	} else if b, ok := e.(mongo.BulkWriteException); ok {
		// insertMany时出现
		for _, we := range b.WriteErrors {
			wes = append(wes, we.WriteError)
		}
	} else {
		return
	}

	uk := _config.UniqueField
	for _, we := range wes {
		m := datas[we.Index].(bson.M)
		id := m[uk]
		if we.Code != 11000 {
			slog.Error("同步写入错误,且非E11000", "错误", e.Error(), "id", id)
			continue
		}
		// 数据存在时, 针对指定字段比较, 如果不一致就修改指定的字段
		if len(_config.CompareField) > 0 {
			var target bson.M
			sinkDB.FindOne(context.Background(), bson.M{uk: id}).Decode(&target)
			update := make(bson.D, 0)
			for _, f := range _config.CompareField {
				origina := m[f]
				if origina != target[f] {
					update = append(update, bson.E{Key: f, Value: origina})
				}
			}
			if len(update) > 0 {
				sinkDB.UpdateOne(context.Background(), bson.M{uk: id}, bson.D{{Key: "$set", Value: update}})
			}
		}
	}

}
func postMoveSuccess(dataIds []any, sourceDb *mongo.Collection, _config config) {
	uk := _config.UniqueField
	for _, dataId := range dataIds {
		switch _config.OnSuccess {
		case "delete":
			sourceDb.DeleteOne(context.TODO(), bson.M{uk: dataId})
		case "markDel":
			sourceDb.UpdateOne(context.TODO(), bson.M{uk: dataId}, bson.D{{Key: "$set", Value: bson.D{{Key: "markDel", Value: 1}}}})
		}
	}
}
