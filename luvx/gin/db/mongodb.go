package db

import (
	"luvx/gin/config"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/func_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	MongoMain = func_x.Lazy(func() *mongo.Database {
		_config := config.AppConfig.MongoDB
		return createMoncgoCli(_config.Uri, _config.Database)
	})
	MongoSlave = func_x.Lazy(func() *mongo.Database {
		_config := config.AppConfig.MongoDB
		return createMoncgoCli(config.Viper.GetString(config.RemoteMongoUri), _config.Database)
	})
	mongoMainMap, mongoSlaveMap = make(map[string]*mongo.Collection, 4), make(map[string]*mongo.Collection, 4)
)

func createMoncgoCli(uri, db string) *mongo.Database {
	defer common_x.TrackTime("初始化MongoDB连接...")()
	remoteClient, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return remoteClient.Database(db)
}

func getCollection(cli *mongo.Database, name string) *mongo.Collection {
	m := common_x.IfThen(cli == MongoSlave.Get(), mongoSlaveMap, mongoMainMap)
	return maps_x.ComputeIfAbsent(m, name, func(name string) *mongo.Collection { return cli.Collection(name) })
}

func GetMainCollection(name string) *mongo.Collection  { return getCollection(MongoMain.Get(), name) }
func GetSlaveCollection(name string) *mongo.Collection { return getCollection(MongoSlave.Get(), name) }
func GetCollectionByName(name string) *mongo.Collection {
	if r, ok := mongoMainMap[name]; ok {
		return r
	}
	if r, ok := mongoSlaveMap[name]; ok {
		return r
	}
	return GetMainCollection(name)
}
