package db

import (
	"context"
	"luvx/gin/config"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoMain, MongoSlave       *mongo.Database
	mongoMainMap, mongoSlaveMap = make(map[string]*mongo.Collection, 4), make(map[string]*mongo.Collection, 4)
)

func init() {
	defer common_x.TrackTime("初始化MongoDB连接...")()
	_config := config.AppConfig.MongoDB
	clientOptions := options.Client().ApplyURI(_config.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	MongoMain = client.Database(_config.Database)

	clientOptions = options.Client().ApplyURI(config.Viper.GetString(config.RemoteMongoUri))
	remoteClient, err := mongo.Connect(context.TODO(), clientOptions)
	MongoSlave = remoteClient.Database(_config.Database)
	if err != nil {
		panic(err)
	}
}

func getCollection(cli *mongo.Database, name string) *mongo.Collection {
	m := common_x.IfThen(cli == MongoSlave, mongoSlaveMap, mongoMainMap)
	return maps_x.ComputeIfAbsent(m, name, func(name string) *mongo.Collection { return cli.Collection(name) })
}

func GetMainCollection(name string) *mongo.Collection  { return getCollection(MongoMain, name) }
func GetSlaveCollection(name string) *mongo.Collection { return getCollection(MongoSlave, name) }
func GetCollectionByName(name string) *mongo.Collection {
	if r, ok := mongoMainMap[name]; ok {
		return r
	}
	if r, ok := mongoSlaveMap[name]; ok {
		return r
	}
	return GetMainCollection(name)
}
