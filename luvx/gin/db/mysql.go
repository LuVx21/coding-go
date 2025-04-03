package db

import (
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/dbs"

	//omysql "github.com/go-sql-driver/mysql"
	"luvx/gin/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var MySQLClient *gorm.DB

func init() {
	defer common_x.TrackTime("初始化MySQL连接...")()

	c := config.AppConfig.MySQL
	//_ = omysql.RegisterTLSConfig("tidb", &tls.Config{
	//    MinVersion: tls.VersionTLS12,
	//    ServerName: c.Host,
	//})
	//format += "&tls=tidb"
	dsn := dbs.MySQLConnectWithDefaultArgs(c.Host, c.Port, c.Username, c.Password, c.Dbname)

	opts := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "t_", // 表名前缀
			SingularTable: true, // 使用单数表名
		},
	}

	var err error
	MySQLClient, err = gorm.Open(mysql.New(mysql.Config{DSN: dsn}), opts)
	if err != nil {
		panic(err)
	}
}
