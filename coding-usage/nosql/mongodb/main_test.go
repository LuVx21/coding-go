package main

import (
    "context"
    "fmt"
    "github.com/luvx21/coding-go/infra/nosql/mongodb"
    where "github.com/pywee/gobson-where"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "testing"
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

    filter := bson.D{bson.E{Key: "age", Value: 0}}
    filter = append(filter, bson.E{Key: "password", Value: "bar_0"})

    opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(100)
    rowsMap, _ := mongodb.RowsMap(context.Background(), collection, filter, opts)
    fmt.Println(rowsMap)
}

func Test_update(t *testing.T) {
    defer beforeAfter("Test_update")()

    filter := bson.D{bson.E{Key: "_id", Value: 0}}
    update := bson.D{{"$set",
        bson.D{{"age", -1}},
    }}
    many, err := collection.UpdateMany(context.TODO(), filter, update)
    fmt.Println(many, err)
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
    distinctValues, _ := collection.Distinct(context.TODO(), "_class", bson.M{"age": 1})
    for _, value := range distinctValues {
        fmt.Println(value)
    }
}

func Test_sort(t *testing.T) {
    defer beforeAfter("Test_sort")()

    opts := options.Find().
        SetProjection(bson.D{{"_id", 1}}).
        SetSort(bson.D{{Key: "age", Value: -1}, {Key: "_id", Value: -1}}).
        SetLimit(100)
    cursor, _ := collection.Find(context.Background(), bson.M{}, opts)
    print(cursor)
}

func Test_01(t *testing.T) {
    opt := where.Parse(`sku!=123 AND (name=456 OR id=789) AND id!=1 ORDER BY name DESC LIMIT 0,10`)
    fmt.Println(opt.Filter)
    fmt.Println(opt.Options)
}
