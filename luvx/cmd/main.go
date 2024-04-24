package main

import (
    "fmt"
    "gorm.io/driver/mysql"
    "gorm.io/gen"
    "gorm.io/gorm"
)

const (
    host, port     = "luvx", 53306
    user, password = "", ""
    dbname         = ""
    url1           = "%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local"
)

type Querier interface {
    FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
    g := gen.NewGenerator(gen.Config{
        OutPath: "./gin/query",
        Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
    })

    db, _ := gorm.Open(mysql.Open(fmt.Sprintf(url1, host, port, user, password, dbname)))
    g.UseDB(db)

    t := g.GenerateModel("user")
    t1 := g.GenerateModel("common_key_value")
    g.ApplyBasic(t, t1)

    //g.ApplyInterface(func(Querier) {}, model.User{}, model.CommonKeyValue{})

    g.Execute()
}
