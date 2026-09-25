package db_mongo

import (
	"log/slog"
	"luvx/gin/config"
	"os"
	"sync"

	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/func_x"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	MongoMainCon  = sync.OnceValue(func() *mongo.Client { return createMoncgoCli(config.AppConfig.MongoDB.Uri) })
	MongoSlaveCon = sync.OnceValue(func() *mongo.Client { return createMoncgoCli(config.Viper.GetString(config.RemoteMongoUri)) })

	MgMainDB  = func_x.Lazy(func() *mongo.Database { return MongoMainCon().Database(config.AppConfig.MongoDB.Database) })
	MgSlaveDB = func_x.Lazy(func() *mongo.Database { return MongoSlaveCon().Database(config.AppConfig.MongoDB.Database) })

	mongoMainMap, mongoSlaveMap = make(map[string]*mongo.Collection, 4), make(map[string]*mongo.Collection, 4)
)

func createMoncgoCli(uri string) *mongo.Client {
	defer common_x.TrackTime("初始化MongoDB连接...")()
	cli, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		slog.Error("Mongo连接失败")
		os.Exit(1)
	}
	return cli
}

func getCollection(cli *mongo.Database, name string) *mongo.Collection {
	m := common_x.IfThen(cli == MgSlaveDB.Get(), mongoSlaveMap, mongoMainMap)
	return maps_x.ComputeIfAbsent(m, name, func(name string) *mongo.Collection { return cli.Collection(name) })
}

func GetMainCollection(name string) *mongo.Collection { return getCollection(MgMainDB.Get(), name) }
func GetSlaveCollection(name string) *mongo.Collection {
	return getCollection(MgSlaveDB.Get(), name)
}
func GetCollectionByName(name string) *mongo.Collection {
	if r, ok := mongoMainMap[name]; ok {
		return r
	}
	if r, ok := mongoSlaveMap[name]; ok {
		return r
	}
	return GetMainCollection(name)
}
