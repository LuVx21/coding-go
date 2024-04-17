package service

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/luvx21/coding-go/coding-common/iterators"
    "github.com/luvx21/coding-go/coding-common/maps_x"
    "github.com/spf13/cast"
    "github.com/valyala/fasthttp"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/exp/slices"
    "luvx/gin/db"
    "luvx/gin/service/common_kv"
    "time"
)

const (
    url = "https://api.bilibili.com/x/space/fav/season/list?season_id=%d&pn=%d&ps=%d"
)

func PullAll() {
    m := common_kv.Get(common_kv.MAP, "bili_season")
    v := m["bili_season"]
    ff := make(map[string]bool)
    _ = json.Unmarshal([]byte(v.CommonValue), &ff)
    for seasonId, flag := range ff {
        if !flag {
            continue
        }
        PullSeasonList(cast.ToInt64(seasonId))
    }
}

func PullSeasonList(seasonId int64) {
    cursor, limit := 1, 20
    iterator := iterators.NewCursorIterator(
        cursor,
        false,
        func(_cursor int) []interface{} {
            return funcName(seasonId, _cursor, limit)
        },
        func(items []interface{}) int {
            if len(items) < limit {
                return -1
            }
            //接口有问题,没有分页取
            //cursor++
            //return cursor
            return -1
        },
        func(i int) bool {
            return i <= 0
        },
    )

    aa := []string{"_id", "title", "pubtime", "bvid", "upper"}
    client := db.MongoDatabase.Collection("bili_video")
    iterator.ForEachRemaining(func(item interface{}) {
        media := item.(map[string]interface{})
        id := media["id"]
        filter := bson.D{{"_id", id}}
        var result bson.M
        _ = client.FindOne(context.TODO(), filter).Decode(&result)
        if result != nil {
            return
        }

        media["_id"] = id
        maps_x.ComputeIfPresent(media, "pubtime", func(k string, v interface{}) interface{} {
            return time.Unix(cast.ToInt64(v), 0)
        })
        maps_x.RemoveIf(media, func(k string, v interface{}) bool {
            return !slices.Contains(aa, k)
        })
        upper := media["upper"].(map[string]interface{})
        maps_x.ComputeIfPresent(upper, "mid", func(k string, v interface{}) interface{} {
            return cast.ToInt64(v)
        })
        upper["seasonId"] = seasonId
        media["invalid"] = 0
        media["from"] = "go"
        _, _ = client.InsertOne(context.TODO(), media)
    })
}

func funcName(seasonId int64, cursor int, limit int) []interface{} {
    client := &fasthttp.Client{}

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    req.SetRequestURI(fmt.Sprintf(url, seasonId, cursor, limit))
    req.Header.SetMethod(fasthttp.MethodGet)
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    _ = client.Do(req, resp)
    ff := make(map[string]interface{})
    _ = json.Unmarshal(resp.Body(), &ff)

    medias := ff["data"].(map[string]interface{})["medias"].([]interface{})
    return medias
}
