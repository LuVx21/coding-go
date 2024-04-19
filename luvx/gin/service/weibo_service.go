package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/gocolly/colly"
    "github.com/luvx21/coding-go/coding-common/ids"
    "github.com/luvx21/coding-go/coding-common/maps_x"
    "github.com/luvx21/coding-go/coding-common/times_x"
    "github.com/spf13/cast"
    "go.mongodb.org/mongo-driver/bson"
    "luvx/gin/db"
)

func PullHotBand() {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL.String())
    })

    c.OnResponse(func(r *colly.Response) {
        ff := make(map[string]interface{})
        _ = json.Unmarshal(r.Body, &ff)

        client := db.MongoDatabase.Collection("weibo_hot_band")
        bandList := ff["data"].(map[string]interface{})["band_list"].([]interface{})
        now := times_x.TimeNowDate()
        worker, _ := ids.NewSnowflakeIdWorker(0, 0)
        for i, v := range bandList {
            vv := v.(map[string]interface{})
            rank := vv["realpos"]
            if rank == nil {
                continue
            }
            word := vv["word"]
            filter := bson.D{{"word", word}}
            var result bson.M
            _ = client.FindOne(context.TODO(), filter).Decode(&result)
            if result != nil {
                rankMap := result["rankMap"].(bson.M)
                oldRank := maps_x.GetOrDefault(rankMap, now, "99")
                if cast.ToInt(oldRank) > cast.ToInt(rank) {
                    rankMap[now] = cast.ToString(rank)
                    update := bson.D{{"$set", bson.D{
                        {"rankMap", rankMap},
                        {"category", vv["category"]},
                    }}}
                    _, _ = client.UpdateOne(context.TODO(), filter, update)
                }
            } else {
                d := bson.D{
                    {"_id", worker.NextId()},
                    {"_class", "org.luvx.boot.tools.dao.mongo.weibo.HotBand"},
                    {"word", word},
                    {"category", vv["category"]},
                    {"rankMap", map[string]string{now: cast.ToString(rank)}},
                }
                _, _ = client.InsertOne(context.TODO(), d)
            }
            fmt.Println(i+1, word, result == nil)
        }
        //fmt.Println("Visited", r.Request.URL.String())
    })

    _ = c.Visit("https://weibo.com/ajax/statuses/hot_band")
}
