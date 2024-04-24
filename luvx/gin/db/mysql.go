package db

import (
    "fmt"
    "github.com/luvx21/coding-go/coding-common/common_x"

    //omysql "github.com/go-sql-driver/mysql"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
    "luvx/gin/config"
)

var format = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"

var MySQLClient *gorm.DB

func init() {
    defer common_x.TrackTime("初始化MySQL连接...")()

    c := config.AppConfig.MySQL
    //_ = omysql.RegisterTLSConfig("tidb", &tls.Config{
    //    MinVersion: tls.VersionTLS12,
    //    ServerName: c.Host,
    //})
    //format += "&tls=tidb"
    dsn := fmt.Sprintf(format, c.Username, c.Password, c.Host, c.Port, c.Dbname)

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
