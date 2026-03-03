package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var db *mongo.Database
var collection *mongo.Collection

func beforeAfter(caseName string) func() {
	if collection == nil {
		_, db, collection = connect()
	}

	return func() {
		fmt.Println(caseName, "test case end...")
	}
}

func print(cursor *mongo.Cursor) {
	defer cursor.Close(context.Background())

	var results []bson.M
	_ = cursor.All(context.Background(), &results)
	for i, result := range results {
		fmt.Println("查询结果:", i, result)
	}
}

func Test_insert(t *testing.T) {
	defer beforeAfter("Test_insert")()

	m := map[string]any{
		"_id":      0,
		"userName": "newnewnew",
	}
	one, err := collection.InsertOne(context.TODO(), m)
	fmt.Println(one.InsertedID, err)
}

func Test_find(t *testing.T) {
	defer beforeAfter("Test_find")()

	day, _ := time.ParseInLocation(time.DateOnly, "2025-05-05", time.UTC)
	day = day.Add(time.Hour * -8)

	filter := bson.D{bson.E{Key: "age", Value: 0}}
	filter = append(filter, bson.E{Key: "name", Value: "/" + "foo" + "/"})
	filter = append(filter, bson.E{Key: "password", Value: "bar_0"})
	filter = append(filter, bson.E{Key: "birthday", Value: bson.M{
		"$gte": day, "$lt": day.AddDate(0, 0, 1),
	}})

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
	rowsMap, _ := mongodb.RowsMap(context.Background(), collection, filter, opts)
	fmt.Println(rowsMap)
}

func Test_update(t *testing.T) {
	defer beforeAfter("Test_update")()

	filter := bson.D{{Key: "_id", Value: 99999}, {Key: "age", Value: 99998}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "age", Value: 18},
			{Key: "invalid", Value: 0},
		}},
	}
	many, err := collection.UpdateMany(context.TODO(), filter, update)
	fmt.Println(many, err)
}
func Test_delete(t *testing.T) {
	defer beforeAfter("Test_delete")()

	collection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: 99999}, {Key: "age", Value: 99998}})
}

func Test_00(t *testing.T) {
	defer beforeAfter("Test_00")()

	var ids = []int64{1, 2, 3}
	filter := bson.D{bson.E{Key: "_id", Value: bson.M{"$in": ids}}}

	//filter := bson.D{{"key1", "value1"}}
	cur, _ := collection.Find(context.TODO(), filter)
	print(cur)

	//for cur.Next(context.Background()) {
	//    var result bson.M
	//    _ = cur.Decode(&result)
	//    fmt.Println(result)
	//}
}

func Test_distinct(t *testing.T) {
	defer beforeAfter("Test_distinct")()
	distinctValues := collection.Distinct(context.TODO(), "_class", bson.M{"age": 1})
	arr, _ := distinctValues.Raw()
	fmt.Println(string(arr))
}

func Test_sort(t *testing.T) {
	defer beforeAfter("Test_sort")()

	opts := options.Find().
		SetProjection(bson.D{{Key: "_id", Value: 1}}).
		SetSort(bson.D{{Key: "age", Value: -1}, {Key: "_id", Value: -1}}).
		SetLimit(100)
	cursor, _ := collection.Find(context.Background(), bson.M{}, opts)
	print(cursor)
}

func Test_upsert_00(t *testing.T) {
	defer beforeAfter("Test_upsert_00")()
	r, e := collection.UpdateOne(context.TODO(), bson.M{"_id": "test1"},
		bson.M{
			"$set": bson.M{
				"expireAt": time.Now().Add(times_x.Day).Unix(),
				"ids":      []string{"a", "b"}},
			"$setOnInsert": bson.M{"createdAt": time.Now().Unix()},
		},
		options.UpdateOne().SetUpsert(true))
	fmt.Println(r, e)

	rr := collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": "test2"},
		bson.M{
			"$set": bson.M{
				"expireAt": time.Now().Add(times_x.Day).Unix(),
				"ids":      []string{"a", "b"}},
			"$inc":         bson.M{"count": 1},
			"$setOnInsert": bson.M{"createdAt": time.Now().Unix()},
		}, options.FindOneAndUpdate().SetUpsert(true),
	)
	fmt.Println(rr.Err())

	// 整体更新, 类似先删除后添加
	_, rrr := collection.ReplaceOne(context.TODO(), bson.M{"_id": "test3"},
		bson.M{
			"expireAt": time.Now().Add(times_x.Day).Unix(),
			"ids":      []string{"a", "b"},
		}, options.Replace().SetUpsert(true),
	)
	fmt.Println(rrr)
}

func Test_count_00(t *testing.T) {
	defer beforeAfter("Test_count_00")()

	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	var m bson.M
	err := db.Collection("config").FindOne(context.TODO(), bson.M{"_id": "app_switch"}, opts).Decode(&m)
	fmt.Println(err == nil)

	c, _ := db.Collection("config").CountDocuments(context.TODO(), bson.M{"_id": "app_switch"})
	fmt.Println(c)
}

func Test_index_00(t *testing.T) {
	defer beforeAfter("Test_index_00")()
	ctx := context.Background()

	indexView := db.Collection("bili_video").Indexes()
	cursor, err := indexView.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	indexCount := 0

	// 遍历索引
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		// 解析索引信息
		indexCount++
		fmt.Printf("\n索引 #%d:\n", indexCount)

		// 显示索引详细信息
		for key, value := range result {
			switch key {
			case "v":
				fmt.Printf("  索引版本: %v\n", value)
			case "key":
				fmt.Printf("  键: %v\n", value)
			case "name":
				fmt.Printf("  索引名称: %v\n", value)
			case "ns":
				fmt.Printf("  命名空间: %v\n", value)
			case "background":
				if value == true {
					fmt.Printf("  后台创建: true\n")
				}
			case "unique":
				if value == true {
					fmt.Printf("  唯一索引: true\n")
				}
			case "sparse":
				if value == true {
					fmt.Printf("  稀疏索引: true\n")
				}
			case "expireAfterSeconds":
				if seconds, ok := value.(int32); ok && seconds > 0 {
					fmt.Printf("  TTL过期时间: %d 秒\n", seconds)
				}
			case "weights":
				fmt.Printf("  权重: %v\n", value)
			case "default_language":
				fmt.Printf("  默认语言: %v\n", value)
			case "language_override":
				fmt.Printf("  语言覆盖: %v\n", value)
			case "textIndexVersion":
				fmt.Printf("  文本索引版本: %v\n", value)
			case "2dsphereIndexVersion":
				fmt.Printf("  2dsphere索引版本: %v\n", value)
			case "bits":
				fmt.Printf("  位: %v\n", value)
			case "min":
				fmt.Printf("  最小值: %v\n", value)
			case "max":
				fmt.Printf("  最大值: %v\n", value)
			case "bucketSize":
				fmt.Printf("  桶大小: %v\n", value)
			default:
				fmt.Printf("  %s: %v\n", key, value)
			}
		}
	}

	if indexCount == 0 {
		fmt.Println("没有找到索引（只有默认的 _id 索引？）")
	} else {
		fmt.Printf("\n总索引数: %d\n", indexCount)
	}
}

func Test_index_01(t *testing.T) {
	defer beforeAfter("Test_index_01")()

	collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]any{"xxx": 1},
		Options: options.Index().SetUnique(true),
	})
}
