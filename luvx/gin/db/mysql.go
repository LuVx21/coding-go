package db

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
    "luvx/config"
)

var MySQLClient *gorm.DB

func init() {
    fmt.Println("初始化数据库连接...")
    var err error
    _config := config.AppConfig.MySQL
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        _config.Username, _config.Password, _config.Host, _config.Dbname)
    MySQLClient, err = gorm.Open(mysql.New(mysql.Config{
        DSN: dsn,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NamingStrategy: schema.NamingStrategy{
            //TablePrefix:   "t_", // 表名前缀
            SingularTable: true, // 使用单数表名
        },
    })

    if err != nil {
        panic(err)
    }
}
