package main

import (
    "context"
    "fmt"
    where "github.com/pywee/gobson-where"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "testing"
)

var db *mongo.Database
var collection *mongo.Collection

func beforeAfter(caseName string) func() {
    if collection == nil {
        _, db, collection = connect()
    }

    return func() {
        fmt.Println(caseName, "end...")
    }
}

func Test_00(t *testing.T) {
    defer beforeAfter("Test_00")()

    var ids = []int64{1, 2, 3}
    filter := bson.D{bson.E{Key: "_id", Value: bson.M{"$in": ids}}}

    //filter := bson.D{{"key1", "value1"}}
    cur, _ := collection.Find(context.TODO(), filter)
    defer cur.Close(context.Background())

    var results []bson.M
    _ = cur.All(context.Background(), &results)
    for i, result := range results {
        fmt.Println("查询结果:", i, result)
    }

    //for cur.Next(context.Background()) {
    //    var result bson.M
    //    _ = cur.Decode(&result)
    //    fmt.Println(result)
    //}
}

func Test_distinct(t *testing.T) {
    defer beforeAfter("Test_distinct")()
    distinctValues, _ := collection.Distinct(context.TODO(), "_class", bson.M{"age": 1})
    for _, value := range distinctValues {
        fmt.Println(value)
    }
}

func Test_01(t *testing.T) {
    opt := where.Parse(`sku!=123 AND (name=456 OR id=789) AND id!=1 ORDER BY name DESC LIMIT 0,10`)
    fmt.Println(opt.Filter)
    fmt.Println(opt.Options)
}
