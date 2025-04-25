package db

import (
	"context"
	"luvx/gin/config"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDatabase *mongo.Database
var collectionMap = make(map[string]*mongo.Collection)

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

func GetCollection(name string) *mongo.Collection {
	return maps_x.ComputeIfAbsent(collectionMap, name, func(name string) *mongo.Collection { return MongoDatabase.Collection(name) })
}
