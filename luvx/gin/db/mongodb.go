package db

import (
    "context"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "luvx/gin/config"
)

var MongoDatabase *mongo.Database

func init() {
    defer common_x.TrackTime("初始化MongoDB连接...")()
    _config := config.AppConfig.MongoDB
    clientOptions := options.Client().ApplyURI(_config.Uri)
    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        panic(err)
    }

    MongoDatabase = client.Database(_config.Database)
}
