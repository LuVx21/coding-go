package main

import (
    "context"
    "fmt"
    "github.com/luvx21/coding-go/coding-usage/nosql"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/event"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

var (
    uri            = nosql.MongoUri
    dbName         = "boot"
    collectionName = "user"
)

func connect() (*mongo.Client, *mongo.Database, *mongo.Collection) {
    clientOptions := options.Client().ApplyURI(uri)

    // 输出查询语句
    var logMonitor = event.CommandMonitor{
        Started: func(ctx context.Context, event *event.CommandStartedEvent) {
            if event.CommandName != "ping" {
                log.Printf("库:%s 命令:%s sql:%+v", event.DatabaseName, event.CommandName, event.Command)
            }
        },
        Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
            if event.CommandName != "ping" {
                log.Printf("查询语句:%s 耗时:%dms", event.CommandName, event.Duration/1000/1000)
            }
        },
        Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
            log.Fatalf("查询语句:%s 耗时:%dms", event.CommandName, event.Duration/1000/1000)
        },
    }
    clientOptions.SetMonitor(&logMonitor)

    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("连接到MongoDB!")

    db := client.Database(dbName)
    collection := db.Collection(collectionName)
    return client, db, collection
}

func main() {
    client, _, collection := connect()
    defer client.Disconnect(context.TODO())

    // 插入文档
    insertResult, _ := collection.InsertOne(context.TODO(), bson.D{
        {"key1", "value1"},
        {"key2", "value2"},
    })
    fmt.Println("插入数据ID:", insertResult.InsertedID)

    // 查询文档
    var result bson.M
    filter := bson.D{{"key1", "value1"}}
    _ = collection.FindOne(context.TODO(), filter).Decode(&result)
    fmt.Println("查询:", result)

    // 更新文档
    update := bson.D{
        {"$set", bson.D{
            {"key2", "updated_value"},
        }},
    }
    updateResult, _ := collection.UpdateOne(context.TODO(), filter, update)
    fmt.Printf("查询到 %v 文档并修改 %v\n", updateResult.MatchedCount, updateResult.ModifiedCount)

    // 删除文档
    deleteResult, _ := collection.DeleteOne(context.TODO(), filter)
    fmt.Printf("删除 %v 文档\n", deleteResult.DeletedCount)
}
