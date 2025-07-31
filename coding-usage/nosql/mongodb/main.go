package main

import (
	"context"
	"fmt"
	"log"

	"github.com/luvx21/coding-go/coding-usage/nosql"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uri            = nosql.MongoUri
	dbName         = "boot"
	collectionName = "user"
)

func connect() (*mongo.Client, *mongo.Database, *mongo.Collection) {
	clientOptions := options.Client().ApplyURI(uri)

	clientOptions.SetMonitor(&mongodb.LogMonitor)

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
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	})
	fmt.Println("插入数据ID:", insertResult.InsertedID)

	// 查询文档
	var result bson.M
	filter := bson.D{{Key: "key1", Value: "value1"}}
	_ = collection.FindOne(context.TODO(), filter).Decode(&result)
	fmt.Println("查询:", result)

	// 更新文档
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "key2", Value: "updated_value"},
		}},
	}
	updateResult, _ := collection.UpdateOne(context.TODO(), filter, update)
	fmt.Printf("查询到 %v 文档并修改 %v\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// 删除文档
	deleteResult, _ := collection.DeleteOne(context.TODO(), filter)
	fmt.Printf("删除 %v 文档\n", deleteResult.DeletedCount)
}
