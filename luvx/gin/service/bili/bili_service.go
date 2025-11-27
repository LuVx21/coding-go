package bili

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/db"
	"luvx/gin/service"
	"luvx/gin/service/common_kv"
	"luvx/gin/service/cookie"

	"github.com/bytedance/sonic"
	gocache_store "github.com/eko/gocache/lib/v4/store"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/iterators"
	"github.com/luvx21/coding-go/coding-common/jsons"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

const (
	ckv_key_up, ckv_key_season = "bili_up", "bili_season"
)

var (
	mixinKeyEncTab = []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
	fields = []string{"_id", "title", "pubtime", "bvid", "upper", "invalid", "from", "pic"}
	cache  = consts.NewLoadableCache(func(ctx context.Context, key any) ([]byte, []gocache_store.Option, error) {
		fmt.Println("自动加载缓存......", key)
		imgKey, subKey := getWbiKeys()
		return []byte(imgKey + "|" + subKey), nil, nil
	})
	rateLimiter               = rate.NewLimiter(0.5, 1)
	mongoClient, mongoClient2 = db.MongoDatabase.Collection("bili_video"), db.RemoteMongoDatabase.Collection("bili_video")
)

func PullAllSeason() {
	m := common_kv.Get(common_kv_dao.MAP, ckv_key_season)
	v := m[ckv_key_season]
	ff := make(map[string]bool)
	_ = sonic.Unmarshal([]byte(v.CommonValue), &ff)
	for _, id := range getCollections() {
		ff[id] = true
	}
	for seasonId, flag := range ff {
		if !flag {
			continue
		}
		runs.Go(func() {
			PullSeasonList(cast_x.ToInt64(seasonId))
		})
	}
}

func PullSeasonList(seasonId int64) {
	opts := options.Find().
		SetProjection(bson.D{{Key: "_id", Value: 1}}).
		SetSort(bson.D{{Key: "_id", Value: -1}}).
		SetLimit(300)
	rowsMap, _ := mongodb.RowsMap(context.Background(), mongoClient, bson.M{"upper.seasonId": seasonId}, opts)
	ids := slices_x.Transfer(func(m bson.M) int64 { return cast_x.ToInt64(m["_id"]) }, *rowsMap...)

	cursor, limit := 1, 20
	iterator := iterators.NewCursorIteratorSimple(
		cursor,
		false,
		func(_cursor int) []any {
			return requestSeasonVideo(seasonId, _cursor, limit)
		},
		func(curId int, items []any) int {
			if len(items) < limit {
				return -1
			}
			// 接口有问题,没有分页取
			// cursor++
			// return cursor
			return -1
		},
		func(i int) bool {
			return i <= 0
		},
	)

	iterator.ForEachRemaining(func(item any) {
		media := item.(map[string]any)
		id := cast_x.ToInt64(media["id"])
		media["_id"] = id
		if slices.Contains(ids, id) {
			return
		}
		filter := bson.D{bson.E{Key: "_id", Value: id}}
		var result bson.M
		_ = mongoClient.FindOne(context.TODO(), filter).Decode(&result)
		if result != nil {
			return
		}

		media["invalid"] = 0
		maps_x.ComputeIfPresent(media, "pubtime", func(k string, v any) any {
			return time.Unix(cast_x.ToInt64(v), 0)
		})
		media["from"] = "go-season"
		upper := media["upper"].(map[string]any)
		maps_x.ComputeIfPresent(upper, "mid", func(k string, v any) any {
			return cast_x.ToInt64(v)
		})
		upper["seasonId"] = seasonId

		maps_x.RemoveIf(media, func(k string, v any) bool {
			return !slices.Contains(fields, k)
		})
		_, _ = mongoClient.InsertOne(context.TODO(), media)
		_, _ = mongoClient2.InsertOne(context.TODO(), media)
		logs.Log.Infoln(media["pubtime"], media["title"])
	})
}

func requestSeasonVideo(seasonId int64, cursor int, limit int) []any {
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
	err := client.Do(req, resp)
	ff := make(map[string]any)
	_ = sonic.Unmarshal(resp.Body(), &ff)

	if cast_x.ToInt32(ff["code"]) != 0 || ff["data"] == nil {
		logs.Log.Warnln("bili->请求结果非json,cookie可能过期", err)
		return make([]any, 0)
	}

	medias := ff["data"].(map[string]any)["medias"].([]any)
	return medias
}

func PullAllUpVideo() {
	service.RunnerLocker.LockRun("拉取bili_up", time.Minute*10, func() {
		midMap := dumpUpId()
		for mid, flag := range midMap {
			if !flag || len(mid) == 0 {
				continue
			}
			// runs.Go(func() {
			PullUpVideo(cast_x.ToInt64(mid))
			// })
		}
	})
}

func PullUpVideo(mid int64) []string {
	opts := options.FindOne().SetSort(bson.M{"pubtime": -1})
	var latest bson.M
	_ = mongoClient.FindOne(context.TODO(), bson.M{"upper.mid": mid}, opts).Decode(&latest)

	cursor, limit := 1, 42
	iterator := iterators.NewCursorIteratorSimple(
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
		id := cast_x.ToInt64(video["aid"])
		if cast_x.ToInt64(video["created"]) <= cast_x.ToInt64(latest["pubtime"])/1000 {
			return
		}

		video["_id"] = id
		filter := bson.D{{Key: "_id", Value: id}}
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
		// _, _ = mongoClient.InsertOne(context.TODO(), video)
		toSave = append(toSave, video)
		result = append(result, video["bvid"].(string))
		logs.Log.Infoln(video["pubtime"], video["title"])
	})
	for _, s := range slices_x.Partition(toSave, 30) {
		_, _ = mongoClient.InsertMany(context.TODO(), s)
		_, _ = mongoClient2.InsertMany(context.TODO(), s)
	}
	return result
}

func requestUpVideo(mid int64, cursor int, limit int) []any {
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
		// "dm_cover_img_str": "",
		// "dm_img_inter":     "",
		// "dm_img_list":      "",
		// "dm_img_str":       "",
	}
	pUrl, _ := nets_x.UrlAddQuery("https://api.bilibili.com/x/space/wbi/arc/search", m)

	newUrlStr, err := signAndGenerateURL(pUrl)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil
	}

	body := biliRequest(newUrlStr.String(), nil, true)
	if body == "" {
		return nil
	}

	ff, _ := jsons.JsonStringToMap[string, any, map[string]any](body)
	vlist := ff["data"].(map[string]any)["list"].(map[string]any)["vlist"].([]any)
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
	json := biliRequest(targetUrl, nil, false)
	g := gjson.Parse(json)

	a := func(s string) string {
		URL := g.Get(s).String()
		return strings.Split(strings.Split(URL, "/")[len(strings.Split(URL, "/"))-1], ".")[0]
	}

	imgKey, subKey := a("data.wbi_img.img_url"), a("data.wbi_img.sub_url")
	return imgKey, subKey
}

func dumpUpId() map[string]bool {
	m := common_kv.Get(common_kv_dao.MAP, ckv_key_up)
	v := m[ckv_key_up]
	midMap := make(map[string]bool)
	_ = sonic.Unmarshal([]byte(v.CommonValue), &midMap)

	mids := append(getFollows(43510), getFollows(-10)...)
	for _, mid := range mids {
		midMap[mid] = true
	}
	return midMap
}

func getFollows(tagid int64) []string {
	tagidStr := cast_x.ToString(tagid)
	v := common_kv.GetOne(common_kv_dao.MAP, "bili_follow")
	g := gjson.Get(v.CommonValue, tagidStr)
	expired := time.Now().Unix() > g.Get("expireAt").Int()
	if !expired {
		return slices_x.Transfer(func(r gjson.Result) string { return r.String() }, g.Get("ids").Array()...)
	}

	body := biliRequest("https://api.bilibili.com/x/relation/tag", map[string]any{
		"tagid":        tagid,
		"pn":           1,
		"ps":           200,
		"mid":          "4489143",
		"web_location": "333.1387",
	}, true)

	mids := gjson.Get(body, "data.#.mid").Array()
	array := make([]string, 0, len(mids))
	for _, mid := range mids {
		array = append(array, mid.Raw)
	}

	common_kv_dao.UpdateJsonMap(common_kv_dao.MAP, "bili_follow",
		"JSON_SET(common_value, ?, CAST(? AS JSON))",
		`$."`+tagidStr+`"`, jsons.ToJsonString(map[string]any{
			"expireAt": time.Now().Add(times_x.Day).Unix(),
			"ids":      array,
		}),
	)

	fmt.Println("哈哈哈", array)
	return array
}

func getCollections() []string {
	cli := db.MongoDatabase.Collection("config")
	var result bson.M
	cli.FindOne(context.TODO(), bson.M{"_id": "bili_season"}).Decode(&result)
	expired := time.Now().Unix() > cast_x.ToInt64(result["expireAt"])
	if !expired {
		return slices_x.Transfer(func(r any) string { return cast_x.ToString(r) }, result["ids"].(bson.A)...)
	}

	ids := make([]string, 0)
	hasMore, pn := true, 1
	for hasMore {
		_json := biliRequest("https://api.bilibili.com/x/v3/fav/folder/collected/list", map[string]any{
			"pn":           pn,
			"ps":           50,
			"up_mid":       "4489143",
			"platform":     "web",
			"web_location": "333.1387",
		}, true)

		for _, id := range gjson.Get(_json, "data.list.#.id").Array() {
			ids = append(ids, id.Raw)
		}
		hasMore = gjson.Get(_json, "data.has_more").Bool()
		pn++
	}

	cli.UpdateOne(context.TODO(), bson.M{"_id": "bili_season"}, bson.M{"$set": bson.M{
		"expireAt": time.Now().Add(times_x.Day).Unix(), "ids": ids,
	}}, options.Update().SetUpsert(true))

	return ids
}

func biliRequest(_url string, queryMap map[string]any, useCookie bool) string {
	pUrl, _ := nets_x.UrlAddQuery(_url, queryMap)

	_ = rateLimiter.Wait(context.TODO())
	logs.Log.Infoln("请求:", pUrl)
	sa := consts.GoRequest.Get(pUrl.String()).
		Set("User-Agent", consts.UserAgent).
		Set("Referer", "https://www.bilibili.com/")
	if useCookie {
		sa = sa.Set("Cookie", cookie.GetCookieStrByHost(".bilibili.com"))
	}
	r, body, errs := sa.End()

	if !sonic.ValidString(body) {
		logs.Log.Warnln("bili->请求结果非json,cookie可能过期", r == nil, body, errs)
		return ""
	}

	return body
}

func timeFlow() {
	const a = "itemOpusStyle,listOnlyfans,opusBigCover,onlyfansVote,decorationCard,onlyfansAssetsV2,forwardListHidden,ugcDelete,onlyfansQaCard,commentsNewVersion,avatarAutoTheme,sunflowerStyle,eva3CardOpus,eva3CardVideo,eva3CardComment"
	const b = `{"platform":"web","device":"pc","spmid":"333.1365"}`

	opts := options.FindOne().SetSort(bson.M{"pubtime": -1})
	var latest bson.M
	_ = mongoClient.FindOne(context.TODO(), bson.M{}, opts).Decode(&latest)

	cursor, offset := 1, ""
	iterator := iterators.NewCursorIterator(
		cursor, false,
		func(_cursor int) gjson.Result {
			c := map[string]any{
				"timezone_offset":        -480,
				"type":                   "video",
				"platform":               "web",
				"page":                   cursor,
				"features":               a,
				"web_location":           "333.1365",
				"x-bili-device-req-json": b,
			}
			if offset != "" {
				c["offset"] = offset
			}
			_json := biliRequest("https://api.bilibili.com/x/polymer/web-dynamic/v1/feed/all", c, true)
			return gjson.Parse(_json)
		},
		func(curId int, items gjson.Result) int {
			array := items.Get("data.items").Array()
			if latest == nil || len(array) == 0 {
				return math.MinInt
			}
			lastTime := array[len(array)-1].Get("modules.module_author.pub_ts").Num
			if cast_x.ToInt64(lastTime) <= cast_x.ToInt64(latest["pubtime"])/1000 {
				return math.MinInt
			}
			hasMore := items.Get("data.has_more").Bool()
			offset = items.Get("data.offset").String()
			cursor++
			return common_x.IfThen(hasMore && cursor <= 30, cursor, math.MinInt)
		},
		func(r gjson.Result) []gjson.Result {
			return r.Get("data.items").Array()
		},
		func(i int) bool {
			return i <= 0
		},
	)

	upIds := dumpUpId()
	toSave := make([]any, 0)
	iterator.ForEachRemaining(func(g gjson.Result) {
		author := g.Get("modules.module_author")
		_, ok := upIds[author.Get("mid").String()]
		if !ok {
			return
		}

		archive := g.Get("modules.module_dynamic.major.archive")
		video := map[string]any{
			"_id":     archive.Get("aid").Int(),
			"title":   archive.Get("title").String(),
			"pubtime": time.Unix(cast_x.ToInt64(author.Get("pub_ts").Num), 0),
			"bvid":    archive.Get("bvid").String(),
			"upper": map[string]any{
				"mid":  cast_x.ToInt64(author.Get("mid").Num),
				"name": author.Get("name").String(),
			},
			"invalid": 0,
			"from":    "dynamic",
			"pic":     archive.Get("cover").String(),
		}

		toSave = append(toSave, video)
	})
	for _, s := range slices_x.Partition(toSave, 30) {
		_, _ = mongoClient.InsertMany(context.TODO(), s)
		_, _ = mongoClient2.InsertMany(context.TODO(), s)
	}
}
