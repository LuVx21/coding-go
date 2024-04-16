package controller

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/luvx21/coding-go/coding-common/common"
    "go.mongodb.org/mongo-driver/bson"
    "luvx/gin/db"
    "luvx/gin/model"
    "luvx/gin/service"
    "net/http"
    "sync"
)

var (
    redisClient = db.InitRedisClient()
)

func HealthyCheck(c *gin.Context) {
    args := 1

    f0 := func() model.User {
        mysql := common.RunWithTime("mysql", func() model.User {
            var user model.User
            if err := db.MySQLClient.Where("id = ?", args).First(&user).Error; err != nil {
                panic(err)
            }
            return user
        })
        return mysql
    }

    f1 := func() bson.M {
        mongo, _ := common.RunWithTime2("mongodb", func() (bson.M, error) {
            userTable := db.MongoDatabase.Collection("user")
            filter := bson.D{{"_id", args}}
            var result bson.M
            a := userTable.FindOne(context.TODO(), filter).Decode(&result)
            return result, a
        })
        return mongo
    }
    f2 := func() string {
        redis, _ := common.RunWithTime2("redis", func() (string, error) {
            return redisClient.Get(context.Background(), "foo").Result()
        })
        return redis
    }

    wg := sync.WaitGroup{}
    r0 := make(chan model.User, 1)
    r1 := make(chan bson.M, 1)
    r2 := make(chan string, 1)
    common.RunInRoutine(&wg, func() { r0 <- f0() })
    common.RunInRoutine(&wg, func() { r1 <- f1() })
    common.RunInRoutine(&wg, func() { r2 <- f2() })

    sqlite, _ := common.RunWithTime2("sqlite", func() ([]map[string]interface{}, error) {
        return db.QueryForMap(db.SqliteClient, "select * from user where id = ?", args)
    })
    cookie := common.RunWithTime("cookie", func() map[string]string {
        return service.GetCookieByHost(".weibo.com", "weibo.com")
    })

    wg.Wait()
    close(r0)
    close(r1)
    close(r2)
    mysql := <-r0
    mongo := <-r1
    redis := <-r2
    c.JSON(http.StatusOK, gin.H{
        "mysql":  mysql,
        "mongo":  mongo,
        "redis":  redis,
        "sqlite": sqlite,
        "cookie": cookie,
    })
}
