package main

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "log"
)

const (
    uri        = ""
    database   = "boot"
    collection = "user"
)

func main() {
    // 设置MongoDB连接信息
    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // 检查是否能够连接MongoDB
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("连接到MongoDB!")

    // 获取数据库和集合的句柄
    database := client.Database(database)
    collection := database.Collection(collection)

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

    // 断开与MongoDB的连接
    err = client.Disconnect(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("连接关闭")
}
