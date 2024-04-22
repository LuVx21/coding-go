package weibo_p

import (
    "context"
    "encoding/json"
    "github.com/gocolly/colly"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/ids"
    "github.com/luvx21/coding-go/coding-common/iterators"
    "github.com/luvx21/coding-go/coding-common/jsons"
    "github.com/luvx21/coding-go/coding-common/logs"
    "github.com/luvx21/coding-go/coding-common/maps_x"
    "github.com/luvx21/coding-go/coding-common/nets_x"
    "github.com/luvx21/coding-go/coding-common/slices_x"
    "github.com/luvx21/coding-go/coding-common/times_x"
    "github.com/tidwall/gjson"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/exp/slices"
    "luvx/gin/common/consts"
    "luvx/gin/db"
    "luvx/gin/service/common_kv"
    "luvx/gin/service/cookie"
    "time"
)

var (
    fields = []string{"_id", "user", "mblogid", "created_at", "text_raw", "text", "retweeted_status", "pic_ids",
        "invalid", "extra",
    }
    collection = db.MongoDatabase.Collection("weibo_feed")
)

func PullHotBand() {
    c := colly.NewCollector()

    c.OnRequest(func(r *colly.Request) {
        logs.Log.Infoln("Visiting", r.URL.String())
    })

    c.OnResponse(func(r *colly.Response) {
        ff := make(map[string]interface{})
        _ = json.Unmarshal(r.Body, &ff)

        client := db.MongoDatabase.Collection("weibo_hot_band")
        bandList := ff["data"].(map[string]interface{})["band_list"].([]interface{})
        now := times_x.TimeNowDate()
        worker, _ := ids.NewSnowflakeIdWorker(0, 0)
        for _, v := range bandList {
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
                if cast_x.ToInt(oldRank) > cast_x.ToInt(rank) {
                    rankMap[now] = cast_x.ToString(rank)
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
                    {"rankMap", map[string]string{now: cast_x.ToString(rank)}},
                }
                _, _ = client.InsertOne(context.TODO(), d)
            }
            //fmt.Println(i+1, word, result == nil)
        }
        //fmt.Println("Visited", r.Request.URL.String())
    })

    _ = c.Visit("https://weibo.com/ajax/statuses/hot_band")
}

func PullByUserAll() {
    key := "weibo_user"
    m := common_kv.Get(common_kv.MAP, key)
    v := m[key]
    ff := make(map[string]bool)
    _ = json.Unmarshal([]byte(v.CommonValue), &ff)
    for uid, flag := range ff {
        if !flag {
            continue
        }
        PullByUser(cast_x.ToInt64(uid))
    }
}

func PullByUser(uid int64) {
    opts := options.FindOne().SetSort(bson.M{"_id": -1})
    var latest bson.M
    _ = collection.FindOne(context.TODO(), bson.M{"user.id": uid}, opts).Decode(&latest)

    cursor := 1
    iterator := iterators.NewCursorIteratorSimple[interface{}, int](
        cursor,
        false,
        func(_cursor int) []interface{} {
            return requestPageOfUser(uid, _cursor)
        },
        func(curId int, items []interface{}) int {
            if latest == nil || items == nil || len(items) == 0 {
                return -1
            }
            id := cast_x.ToInt64(items[len(items)-1].(map[string]interface{})["id"])
            if id <= latest["_id"].(int64) {
                return -1
            }
            cursor++
            // 最大100页,避免一直取
            return common_x.IfThen(cursor <= 100, cursor, -1)
        },
        func(i int) bool {
            return i <= 0
        },
    )
    iterator.ForEachRemaining(func(item interface{}) {
        feed := item.(map[string]interface{})
        id := cast_x.ToInt64(feed["id"])
        if latest["_id"] != nil && id <= latest["_id"].(int64) {
            return
        }
        // ---------------------------------------------

        ret := feed["retweeted_status"]
        if ret != nil {
            feed["retweeted_status"] = parseAndSaveFeed(ret.(map[string]interface{}), true)
        }
        parseAndSaveFeed(feed, false)
    })
}

func PullByGroup() {
    var groupId int64 = 4670120389774996
    opts := options.FindOne().SetSort(bson.M{"_id": -1})
    var latest bson.M
    _ = collection.FindOne(context.TODO(), bson.M{}, opts).Decode(&latest)

    var cursor int64 = 0
    iterator := iterators.NewCursorIterator[interface{}, int64, common_x.Pair[[]interface{}, int64]](
        cursor, false,
        func(_cursor int64) common_x.Pair[[]interface{}, int64] {
            return requestPageOfGroup(groupId, _cursor)
        },
        func(curId int64, p common_x.Pair[[]interface{}, int64]) int64 {
            items := p.K
            if latest == nil || items == nil || len(items) == 0 {
                return -1
            }
            last := items[len(items)-1]
            id := cast_x.ToInt64(last.(map[string]interface{})["id"])
            if id <= latest["_id"].(int64) {
                return -1
            }
            return p.V
        },
        func(p common_x.Pair[[]interface{}, int64]) []interface{} {
            return p.K
        },
        func(i int64) bool {
            return i <= 0
        },
    )
    iterator.ForEachRemaining(func(item interface{}) {
        feed := item.(map[string]interface{})
        id := cast_x.ToInt64(feed["id"])
        if latest["_id"] != nil && id <= latest["_id"].(int64) {
            return
        }
        ret := feed["retweeted_status"]
        if ret != nil {
            feed["retweeted_status"] = parseAndSaveFeed(ret.(map[string]interface{}), true)
        }
        parseAndSaveFeed(feed, false)
    })
}

func parseAndSaveFeed(feed map[string]interface{}, retweeted bool) int64 {
    id := cast_x.ToInt64(feed["id"])
    feed["_id"] = id
    var r bson.M
    _ = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&r)
    if r != nil {
        return id
    }

    //feed["extra"] = map[string]interface{}{
    //    "retweeted": retweeted,
    //}

    if cast_x.ToBool(feed["isLongText"]) {
        feed["text"] = requestLongText(feed["mblogid"].(string))
    }

    picUrl := make([]string, 0)
    i2 := feed["pic_ids"]
    if i2 != nil {
        if b, picIds := slices_x.IsEmpty[[]interface{}, interface{}](i2.([]interface{})); !b {
            i := feed["pic_infos"]
            if i != nil {
                picInfos := i.(map[string]interface{})
                for _, picId := range picIds {
                    url := picInfos[picId.(string)].(map[string]interface{})["largest"].(map[string]interface{})["url"]
                    picUrl = append(picUrl, url.(string))
                }
            } else {
                i3 := feed["mix_media_info"]
                if i3 != nil {
                    i4 := i3.(map[string]interface{})["items"].([]interface{})
                    for _, i5 := range i4 {
                        m := i5.(map[string]interface{})
                        if m["type"] == "pic" {
                            url := m["data"].(map[string]interface{})["largest"].(map[string]interface{})["url"]
                            picUrl = append(picUrl, url.(string))
                        }
                    }
                }
            }
        }
    }
    feed["pic_ids"] = picUrl
    // 视频: page_info.media_info.{h5_url,playback_list}

    maps_x.ComputeIfPresent(feed, "created_at", func(k string, v interface{}) interface{} {
        t, _ := time.ParseInLocation(time.RubyDate, v.(string), time.Local)
        return t
    })
    i := feed["user"]
    if i != nil {
        user := i.(map[string]interface{})
        feed["user"] = map[string]interface{}{
            "id":   cast_x.ToInt64(user["id"]),
            "name": user["screen_name"],
        }
    }
    //feed["invalid"] = 0
    maps_x.RemoveIf(feed, func(k string, v interface{}) bool {
        return !slices.Contains(fields, k)
    })
    _, _ = collection.InsertOne(context.TODO(), feed)
    return id
}

func requestLongText(mblogid string) string {
    m := map[string]any{
        "id": mblogid,
    }
    pUrl, _ := nets_x.UrlAddQuery("https://weibo.com/ajax/statuses/longtext", m)
    _ = consts.RateLimiter.Wait(context.TODO())
    _, body, _ := consts.GoRequest.Get(pUrl.String()).
        Set("User-Agent", consts.UserAgent).
        Set("Host", "weibo.com").
        Set("Cookie", getCookie()).
        End()
    return gjson.Get(body, "data.longTextContent").String()
}

func requestPageOfUser(uid int64, cursor int) []interface{} {
    m := map[string]any{
        "uid":     uid,
        "page":    cursor,
        "feature": 0,
    }
    pUrl, _ := nets_x.UrlAddQuery("https://weibo.com/ajax/statuses/mymblog", m)
    logs.Log.Infoln("请求:", pUrl)
    _ = consts.RateLimiter.Wait(context.TODO())
    _, body, _ := consts.GoRequest.Get(pUrl.String()).
        Set("User-Agent", consts.UserAgent).
        Set("Host", "weibo.com").
        Set("Cookie", getCookie()).
        End()

    ff, _ := jsons.JsonStringToMap[string, interface{}, map[string]interface{}](body)
    i := ff["data"]
    if i == nil {
        return make([]interface{}, 0)
    }
    list := i.(map[string]interface{})["list"].([]interface{})
    return list
}

func requestPageOfGroup(groupId int64, cursor int64) common_x.Pair[[]interface{}, int64] {
    m := map[string]any{
        "list_id":      groupId,
        "max_id":       cursor,
        "refresh":      4,
        "fast_refresh": 1,
        "count":        25,
    }
    pUrl, _ := nets_x.UrlAddQuery("https://weibo.com/ajax/feed/groupstimeline", m)
    _ = consts.RateLimiter.Wait(context.TODO())
    r, body, _ := consts.GoRequest.Get(pUrl.String()).
        Set("User-Agent", consts.UserAgent).
        Set("Host", "weibo.com").
        Set("Cookie", getCookie()).
        End()
    logs.Log.Infof("请求: %s 响应:%s\n", pUrl, r.Status)

    ff, _ := jsons.JsonStringToMap[string, interface{}, map[string]interface{}](body)
    list := ff["statuses"].([]interface{})
    maxId := cast_x.ToInt64(ff["max_id"])

    return common_x.NewPair[[]interface{}, int64](list, maxId)
}

func getCookie() string {
    return cookie.GetCookieStrByHost(".weibo.com", "weibo.com")
}
