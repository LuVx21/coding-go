package rss

import (
    "context"
    "github.com/luvx21/coding-go/coding-common/logs"
    "github.com/luvx21/coding-go/coding-common/slices_x"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "luvx/gin/db"
    "luvx/gin/service/common_kv"
    "luvx/gin/service/soup"
)

const ()

var (
    collection = db.MongoDatabase.Collection("rss_feed1")
)

func PullByKey() {
    m := common_kv.Get(8)
    for k, v := range m {
        go func() {
            logs.Log.Infoln("spider拉取:", k)
            items := spiderIndexPage(k, v.CommonValue)
            _, _ = collection.InsertMany(context.TODO(), slices_x.ToAnySliceE(items...))
        }()
    }
}

func spiderIndexPage(key, paramJson string) []soup.PageContent {
    filter := bson.M{
        "spiderKey": key,
        "content":   bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
    }
    opts := options.Find().
        SetProjection(bson.D{
            {"url", 1},
            //{"content", 1},
        }).
        SetSort(bson.M{"_id": -1}).SetLimit(2000)
    cursor, _ := collection.Find(context.Background(), filter, opts)
    defer cursor.Close(context.Background())

    ignoreIndexItemUrlSet := make([]string, 0, 2000)
    for cursor.Next(context.Background()) {
        m := make(map[string]interface{})
        _ = cursor.Decode(&m)
        ignoreIndexItemUrlSet = append(ignoreIndexItemUrlSet, m["url"].(string))
    }

    return soup.Of(paramJson).
        SetIgnoreIndexItemUrlSet(ignoreIndexItemUrlSet).
        Visit()
}
