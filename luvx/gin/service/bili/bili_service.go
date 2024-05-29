package bili

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "github.com/bytedance/sonic"
    "github.com/luvx21/coding-go/coding-common/cast_x"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "github.com/luvx21/coding-go/coding-common/iterators"
    "github.com/luvx21/coding-go/coding-common/jsons"
    "github.com/luvx21/coding-go/coding-common/logs"
    "github.com/luvx21/coding-go/coding-common/maps_x"
    "github.com/luvx21/coding-go/coding-common/nets_x"
    "github.com/luvx21/coding-go/coding-common/slices_x"
    "github.com/luvx21/coding-go/infra/nosql/mongodb"
    "github.com/tidwall/gjson"
    "github.com/valyala/fasthttp"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/exp/slices"
    "golang.org/x/time/rate"
    "luvx/gin/common/consts"
    "luvx/gin/db"
    "luvx/gin/service/common_kv"
    "luvx/gin/service/cookie"
    "net/url"
    "sort"
    "strconv"
    "strings"
    "time"
)

var (
    mixinKeyEncTab = []int{
        46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
        33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
        61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
        36, 20, 34, 44, 52,
    }
    fields = []string{"_id", "title", "pubtime", "bvid", "upper", "invalid", "from"}
    cache  = consts.NewLoadableCache[[]byte](func(ctx context.Context, key any) ([]byte, error) {
        fmt.Println("自动加载缓存......", key)
        imgKey, subKey := getWbiKeys()
        return []byte(imgKey + "|" + subKey), nil
    })
    rateLimiter = rate.NewLimiter(1, 1)
    mongoClient = db.MongoDatabase.Collection("bili_video")
)

func PullAll() {
    key := "bili_season"
    m := common_kv.Get(common_kv.MAP, key)
    v := m[key]
    ff := make(map[string]bool)
    _ = sonic.Unmarshal([]byte(v.CommonValue), &ff)
    for seasonId, flag := range ff {
        if !flag {
            continue
        }
        go PullSeasonList(cast_x.ToInt64(seasonId))
    }
}

func PullSeasonList(seasonId int64) {
    opts := options.Find().
        SetProjection(bson.D{{"_id", 1}}).
        SetSort(bson.D{{Key: "_id", Value: -1}}).
        SetLimit(300)
    rowsMap, _ := mongodb.RowsMap(context.Background(), mongoClient, bson.M{"upper.seasonId": seasonId}, opts)
    ids := slices_x.Transfer(func(m bson.M) int64 { return cast_x.ToInt64(m["_id"]) }, *rowsMap...)

    cursor, limit := 1, 20
    iterator := iterators.NewCursorIteratorSimple(
        cursor,
        false,
        func(_cursor int) []interface{} {
            return requestSeasonVideo(seasonId, _cursor, limit)
        },
        func(curId int, items []interface{}) int {
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

    iterator.ForEachRemaining(func(item interface{}) {
        media := item.(map[string]interface{})
        id := cast_x.ToInt64(media["id"])
        media["_id"] = id
        if slices.Contains(ids, id) {
            return
        }
        filter := bson.D{{"_id", id}}
        var result bson.M
        _ = mongoClient.FindOne(context.TODO(), filter).Decode(&result)
        if result != nil {
            return
        }

        media["invalid"] = 0
        maps_x.ComputeIfPresent(media, "pubtime", func(k string, v interface{}) interface{} {
            return time.Unix(cast_x.ToInt64(v), 0)
        })
        media["from"] = "go-season"
        upper := media["upper"].(map[string]interface{})
        maps_x.ComputeIfPresent(upper, "mid", func(k string, v interface{}) interface{} {
            return cast_x.ToInt64(v)
        })
        upper["seasonId"] = seasonId

        maps_x.RemoveIf(media, func(k string, v interface{}) bool {
            return !slices.Contains(fields, k)
        })
        _, _ = mongoClient.InsertOne(context.TODO(), media)
        logs.Log.Infoln(media["pubtime"], media["title"])
    })
}

func requestSeasonVideo(seasonId int64, cursor int, limit int) []interface{} {
    client := &fasthttp.Client{}

    req := fasthttp.AcquireRequest()
    defer fasthttp.ReleaseRequest(req)
    pUrl, _ := nets_x.UrlAddQuery("https://api.bilibili.com/x/space/fav/season/list", map[string]any{
        "season_id": seasonId,
        "pn":        cursor,
        "ps":        limit,
    })
    req.SetRequestURI(pUrl.String())
    req.Header.SetMethod(fasthttp.MethodGet)
    resp := fasthttp.AcquireResponse()
    defer fasthttp.ReleaseResponse(resp)

    _ = rateLimiter.Wait(context.TODO())
    logs.Log.Infoln("请求:", pUrl)
    _ = client.Do(req, resp)
    ff := make(map[string]interface{})
    _ = sonic.Unmarshal(resp.Body(), &ff)

    medias := ff["data"].(map[string]interface{})["medias"].([]interface{})
    return medias
}

func PullAllUpVideo() {
    key := "bili_up"
    m := common_kv.Get(common_kv.MAP, key)
    v := m[key]
    ff := make(map[string]bool)
    _ = sonic.Unmarshal([]byte(v.CommonValue), &ff)
    for mid, flag := range ff {
        if !flag {
            continue
        }
        go PullUpVideo(cast_x.ToInt64(mid))
    }
}

func PullUpVideo(mid int64) []string {
    opts := options.FindOne().SetSort(bson.M{"pubtime": -1})
    var latest bson.M
    _ = mongoClient.FindOne(context.TODO(), bson.M{"upper.mid": mid}, opts).Decode(&latest)

    cursor, limit := 1, 30
    iterator := iterators.NewCursorIteratorSimple[any, int](
        cursor,
        false,
        func(_cursor int) []any {
            return requestUpVideo(mid, _cursor, limit)
        },
        func(curId int, items []any) int {
            if latest == nil || len(items) < limit {
                return -1
            }
            last := cast_x.ToInt64(items[len(items)-1].(map[string]any)["created"])
            if last <= cast_x.ToInt64(latest["pubtime"])/1000 {
                return -1
            }
            cursor++
            return common_x.IfThen(cursor <= 100, cursor, -1)
        },
        func(i int) bool {
            return i <= 0
        },
    )
    result, toSave := make([]string, 0), make([]any, 0)
    iterator.ForEachRemaining(func(item any) {
        video := item.(map[string]any)
        id := video["aid"]
        if cast_x.ToInt64(video["created"]) <= cast_x.ToInt64(latest["pubtime"])/1000 {
            return
        }

        video["_id"] = id
        filter := bson.D{{"_id", id}}
        var r bson.M
        _ = mongoClient.FindOne(context.TODO(), filter).Decode(&r)
        if r != nil {
            return
        }

        video["invalid"] = 0
        video["pubtime"] = time.Unix(cast_x.ToInt64(video["created"]), 0)
        video["from"] = "go-up"
        upper := map[string]any{
            "mid":  cast_x.ToInt64(video["mid"]),
            "name": video["author"],
        }
        if v, ok := video["season_id"]; ok {
            upper["seasonId"] = cast_x.ToInt64(v)
        }
        video["upper"] = upper

        maps_x.RemoveIf(video, func(k string, v any) bool {
            return !slices.Contains(fields, k)
        })
        //_, _ = mongoClient.InsertOne(context.TODO(), video)
        toSave = append(toSave, video)
        result = append(result, video["bvid"].(string))
        logs.Log.Infoln(video["pubtime"], video["title"])
    })
    for _, s := range slices_x.Partition(toSave, 30) {
        _, _ = mongoClient.InsertMany(context.TODO(), s)
    }
    return result
}

func requestUpVideo(mid int64, cursor int, limit int) []interface{} {
    m := map[string]any{
        "mid":           mid,
        "ps":            limit,
        "tid":           "0",
        "special_type":  "",
        "pn":            cursor,
        "keyword":       "",
        "order":         "pubdate",
        "platform":      "web",
        "web_location":  "1550101",
        "order_avoided": "true",
        //"dm_cover_img_str": "",
        //"dm_img_inter":     "",
        //"dm_img_list":      "",
        //"dm_img_str":       "",
    }
    pUrl, _ := nets_x.UrlAddQuery("https://api.bilibili.com/x/space/wbi/arc/search", m)

    newUrlStr, err := signAndGenerateURL(pUrl)
    if err != nil {
        fmt.Printf("Error: %s", err)
        return nil
    }

    cookie := cookie.GetCookieStrByHost(".bilibili.com")
    _ = rateLimiter.Wait(context.TODO())
    logs.Log.Infoln("请求:", pUrl)
    _, body, _ := consts.GoRequest.Get(newUrlStr.String()).
        Set("User-Agent", consts.UserAgent).
        Set("Referer", "https://www.bilibili.com/").
        Set("Cookie", cookie).
        End()

    ff, _ := jsons.JsonStringToMap[string, interface{}, map[string]interface{}](body)
    vlist := ff["data"].(map[string]interface{})["list"].(map[string]interface{})["vlist"].([]interface{})
    return vlist
}

func signAndGenerateURL(pUrl *url.URL) (*url.URL, error) {
    query := pUrl.Query()
    params := map[string]string{}
    for k, v := range query {
        params[k] = v[0]
    }
    imgKey, subKey := getWbiKeysCached()
    newParams := encWbi(params, imgKey, subKey)
    for k, v := range newParams {
        query.Set(k, v)
    }
    pUrl.RawQuery = query.Encode()
    return pUrl, nil
}

// 来源: https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/user/space.md#%E6%9F%A5%E8%AF%A2%E7%94%A8%E6%88%B7%E6%8A%95%E7%A8%BF%E8%A7%86%E9%A2%91%E6%98%8E%E7%BB%86
func encWbi(params map[string]string, imgKey, subKey string) map[string]string {
    mixinKey := getMixinKey(imgKey + subKey)
    currTime := strconv.FormatInt(time.Now().Unix(), 10)
    params["wts"] = currTime

    // Sort keys
    keys := make([]string, 0, len(params))
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    // Remove unwanted characters
    for k, v := range params {
        v = sanitizeString(v)
        params[k] = v
    }

    // Build URL parameters
    query := url.Values{}
    for _, k := range keys {
        query.Set(k, params[k])
    }
    queryStr := query.Encode()

    // Calculate w_rid
    hash := md5.Sum([]byte(queryStr + mixinKey))
    params["w_rid"] = hex.EncodeToString(hash[:])
    return params
}

func getMixinKey(orig string) string {
    var str strings.Builder
    for _, v := range mixinKeyEncTab {
        if v < len(orig) {
            str.WriteByte(orig[v])
        }
    }
    return str.String()[:32]
}

func sanitizeString(s string) string {
    unwantedChars := []string{"!", "'", "(", ")", "*"}
    for _, char := range unwantedChars {
        s = strings.ReplaceAll(s, char, "")
    }
    return s
}

func getWbiKeysCached() (string, string) {
    sign, _ := cache.Get(context.TODO(), "sign")
    split := strings.Split(string(sign), "|")
    return split[0], split[1]
}

func getWbiKeys() (string, string) {
    targetUrl := "https://api.bilibili.com/x/web-interface/nav"
    logs.Log.Infoln("请求:", targetUrl)
    _, json, _ := consts.GoRequest.Get(targetUrl).
        Set("User-Agent", consts.UserAgent).
        Set("Referer", "https://www.bilibili.com/").
        End()

    imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
    subURL := gjson.Get(json, "data.wbi_img.sub_url").String()
    imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
    subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]
    return imgKey, subKey
}
