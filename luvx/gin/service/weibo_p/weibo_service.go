package weibo_p

import (
	"context"
	"fmt"
	"log/slog"
	"luvx_service_sdk/proto_gen/proto_kv"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/db"
	"luvx/gin/model"
	"luvx/gin/service"
	commonkvservice "luvx/gin/service/common_kv"
	"luvx/gin/service/cookie"
	"luvx/gin/service/rpc"

	"github.com/bytedance/sonic"
	"github.com/gocolly/colly"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	. "github.com/luvx21/coding-go/coding-common/common_x/alias_x"
	. "github.com/luvx21/coding-go/coding-common/common_x/pairs"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/ids"
	"github.com/luvx21/coding-go/coding-common/iterators"
	"github.com/luvx21/coding-go/coding-common/jsons"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/maths_x"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/luvx21/coding-go/infra/logs"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

var (
	fields = []string{"_id", "user", "mblogid", "created_at", "text_raw", "text", "retweeted_status", "pic_ids",
		"invalid", "read", "extra", "groupId", "mix_media_info", "page_info",
	}
	// collection = db.MongoDatabase.Collection("weibo_feed")
	collection = db.GetCollection("weibo_feed")
)

func PullHotBand() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		logs.Log.Infoln("请求:", r.URL.String())
		r.Headers.Add("Referer", "https://weibo.com/hot/search")
	})

	c.OnResponse(func(r *colly.Response) {
		ff := make(map[string]any)
		_ = sonic.Unmarshal(r.Body, &ff)
		fmt.Println("结果:", ff)
		client := db.MongoDatabase.Collection("weibo_hot_band")
		bandList := ff["data"].(map[string]any)["band_list"].([]any)
		now := times_x.TimeNowDate()
		worker, _ := ids.NewSnowflakeIdWorker(0, 0)
		for _, v := range bandList {
			vv := v.(map[string]any)
			rank := vv["realpos"]
			if rank == nil {
				continue
			}
			word := vv["word"]
			filter := bson.D{{Key: "word", Value: word}}
			var result bson.M
			_ = client.FindOne(context.TODO(), filter).Decode(&result)
			if result != nil {
				rankMap := result["rankMap"].(bson.M)
				oldRank := maps_x.GetOrDefault(rankMap, now, "99")
				if cast_x.ToInt(oldRank) > cast_x.ToInt(rank) {
					rankMap[now] = cast_x.ToString(rank)
					update := bson.D{{Key: "$set", Value: bson.D{
						{Key: "rankMap", Value: rankMap},
						{Key: "category", Value: vv["category"]},
					}}}
					_, _ = client.UpdateOne(context.TODO(), filter, update)
				}
			} else {
				d := bson.D{
					{Key: "_id", Value: worker.NextId()},
					// {Key: "_class", Value: "org.luvx.boot.tools.dao.mongo.weibo.HotBand"},
					{Key: "word", Value: word},
					{Key: "category", Value: vv["category"]},
					{Key: "rankMap", Value: map[string]string{now: cast_x.ToString(rank)}},
				}
				_, _ = client.InsertOne(context.TODO(), d)
			}
			// fmt.Println(i+1, word, result == nil)
		}
		// fmt.Println("Visited", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		logs.Log.Errorln("请求异常:", r.Request.URL.RequestURI(), err.Error())
	})

	// https://weibo.com/ajax/side/hotSearch
	_ = c.Visit("https://weibo.com/ajax/statuses/hot_band")
}

func PullByUserAll() {
	key := "weibo_user"
	m := commonkvservice.Get(common_kv_dao.MAP, key)
	v := m[key]
	ff := make(map[string]bool)
	_ = sonic.Unmarshal([]byte(v.CommonValue), &ff)
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
	iterator := iterators.NewCursorIteratorSimple(
		cursor,
		false,
		func(_cursor int) []any {
			return requestPageOfUser(uid, _cursor)
		},
		func(curId int, items []any) int {
			if latest == nil || items == nil || len(items) == 0 {
				return -1
			}
			id := cast_x.ToInt64(items[len(items)-1].(map[string]any)["id"])
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
	arr := make([]any, 0)
	iterator.ForEachRemaining(func(item any) {
		feed := item.(map[string]any)
		id := cast_x.ToInt64(feed["id"])
		if latest["_id"] != nil && id <= latest["_id"].(int64) {
			return
		}
		// ---------------------------------------------

		ret := feed["retweeted_status"]
		if ret != nil {
			f := ret.(map[string]any)
			feed["retweeted_status"] = parseAndSaveFeed(f, true)
			arr = append(arr, f)
		}
		parseAndSaveFeed(feed, false)
		arr = append(arr, feed)
	})
	for i := len(arr) - 1; i >= 0; i-- {
		one, err := collection.InsertOne(context.TODO(), arr[i])
		if err != nil {
			logs.Log.Infoln("weibo PullByUser:", err)
			continue
		}
		logs.Log.Infoln("weibo PullByUser:", one.InsertedID)
	}
}

func PullByGroupLock() {
	service.RunnerLocker.LockRun("拉取分组微博", time.Minute*3, func() {
		PullByGroup(4670120389774996)
		PullByGroup(3639801313908027)
	})
}

func PullByGroup(groupId int64) {
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	var latest bson.M
	_ = collection.FindOne(context.TODO(), bson.M{"groupId": groupId}, opts).Decode(&latest)

	var cursor int64 = 0
	iterator := iterators.NewCursorIterator(
		cursor, false,
		func(_cursor int64) Pair[[]any, int64] {
			return requestPageOfGroup(groupId, _cursor)
		},
		func(curId int64, p Pair[[]any, int64]) int64 {
			items := p.K
			if latest == nil || items == nil || len(items) == 0 {
				return -1
			}
			last := items[len(items)-1]
			id := cast_x.ToInt64(last.(map[string]any)["id"])
			if id <= latest["_id"].(int64) {
				return -1
			}
			return p.V
		},
		func(p Pair[[]any, int64]) []any {
			return p.K
		},
		func(i int64) bool {
			return i <= 0
		},
	)
	feeds := make([]any, 0)
	iterator.ForEachRemaining(func(item any) {
		feed := item.(map[string]any)
		id := cast_x.ToInt64(feed["id"])
		if latest["_id"] != nil && id <= latest["_id"].(int64) {
			return
		}
		ret := feed["retweeted_status"]
		if ret != nil {
			f := ret.(map[string]any)
			feed["retweeted_status"] = parseAndSaveFeed(f, true)
			feeds = append(feeds, f)
		}
		parseAndSaveFeed(feed, false)
		feed["groupId"] = groupId
		feeds = append(feeds, feed)
	})

	runs.Go(func() {
		if rpc.KvRpcClient != nil {
			urls := slices_x.FlatMap(feeds, func(feed any) []string { return feed.(map[string]any)["pic_ids"].([]string) })
			toSave := sets.NewSet(urls...)
			for _url := range *toSave {
				_, err := rpc.KvRpcClient.Get(context.TODO(), &proto_kv.Key{Key: _url})
				if err == nil {
					continue
				}

				_ = consts.RateLimiter.Wait(context.TODO())
				_, body, errors := consts.GoRequest.Get(_url).EndBytes()
				if len(errors) > 0 {
					continue
				}
				_, err = rpc.KvRpcClient.Put(context.TODO(), &proto_kv.PutRequest{Entry: &proto_kv.Entry{Key: _url, Value: body}, Expire: int64(7 * 24 * time.Hour.Seconds())})
				if err != nil {
					slog.Warn("rpc put", "err", err.Error())
				}
			}
		}
	})

	arrs := slices_x.Partition(feeds, 5)
	for i := len(arrs) - 1; i >= 0; i-- {
		arr := arrs[i]
		if many, e := collection.InsertMany(context.TODO(), arr); e != nil {
			for j := len(arr) - 1; j >= 0; j-- {
				if one, err := collection.InsertOne(context.TODO(), arr[j]); err != nil {
					// logs.Log.Infoln("weibo insert1:", err)
					continue
				} else {
					logs.Log.Infoln("weibo insert1:", one.InsertedID)
				}
			}
		} else {
			logs.Log.Infoln("weibo insert", len(arr), len(many.InsertedIDs))
		}
	}
}

func parseAndSaveFeed(feed map[string]any, retweeted bool) int64 {
	id := cast_x.ToInt64(feed["id"])
	feed["_id"] = id
	// var r bson.M
	// _ = collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&r)
	// if r != nil {
	// 	return id
	// }

	// feed["extra"] = map[string]any{
	//    "retweeted": retweeted,
	// }

	if cast_x.ToBool(feed["isLongText"]) {
		feed["text"] = requestLongText(feed["mblogid"].(string))
	}

	picUrl := make([]string, 0)
	i2 := feed["pic_ids"]
	if i2 != nil {
		if b, picIds := slices_x.IsEmpty(i2.([]any)); !b {
			i := feed["pic_infos"]
			if i != nil {
				picInfos := i.(map[string]any)
				for _, picId := range picIds {
					url := picInfos[picId.(string)].(map[string]any)["largest"].(map[string]any)["url"]
					picUrl = append(picUrl, url.(string))
				}
			} else {
				i3 := feed["mix_media_info"]
				if i3 != nil {
					i4 := i3.(map[string]any)["items"].([]any)
					for _, i5 := range i4 {
						m := i5.(map[string]any)
						if m["type"] == "pic" {
							url := m["data"].(map[string]any)["largest"].(map[string]any)["url"]
							picUrl = append(picUrl, url.(string))
						}
					}
				}
			}
		}
	}
	feed["pic_ids"] = picUrl
	// 视频: page_info.media_info.{h5_url,playback_list}

	maps_x.ComputeIfPresent(feed, "created_at", func(k string, v any) any {
		t, _ := time.ParseInLocation(time.RubyDate, v.(string), time.Local)
		return t
	})
	i := feed["user"]
	if i != nil {
		user := i.(map[string]any)
		feed["user"] = map[string]any{
			"id":   cast_x.ToInt64(user["id"]),
			"name": user["screen_name"],
		}
	} else {
		feed["user"] = map[string]any{"id": 0, "name": ""}
	}
	feed["invalid"] = 0
	feed["read"] = 0
	if retweeted {
		feed["invalid"] = 1
	}
	maps_x.RemoveIf(feed, func(k string, v any) bool {
		return !slices.Contains(fields, k)
	})
	// _, _ = collection.InsertOne(context.TODO(), feed)
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
		Set("referer", "https://weibo.com/mygroups?gid=4670120389774996").
		End()
	return gjson.Get(body, "data.longTextContent").String()
}

func requestPageOfUser(uid int64, cursor int) []any {
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
		Set("referer", "https://weibo.com/mygroups?gid=4670120389774996").
		End()

	ff, _ := jsons.JsonStringToMap[string, any, map[string]any](body)
	i := ff["data"]
	if i == nil {
		return make([]any, 0)
	}
	list := i.(map[string]any)["list"].([]any)
	return list
}

func requestPageOfGroup(groupId int64, cursor int64) Pair[[]any, int64] {
	m := map[string]any{
		"list_id":      groupId,
		"max_id":       cursor,
		"refresh":      4,
		"fast_refresh": 1,
		"count":        25,
	}
	pUrl, _ := nets_x.UrlAddQuery("https://weibo.com/ajax/feed/groupstimeline", m)
	_ = consts.RateLimiter.Wait(context.TODO())
	response, body, errors := consts.GoRequest.Get(pUrl.String()).
		Set("User-Agent", consts.UserAgent).
		Set("Host", "weibo.com").
		Set("Cookie", getCookie()).
		Set("referer", "https://weibo.com/mygroups?gid=4670120389774996").
		End()
	isJson := sonic.ValidString(body)
	logs.Log.Infof("请求: %s 响应:%v", pUrl, isJson)
	if !isJson {
		logs.Log.Warnln("weibo->请求结果非json,cookie可能过期", response == nil, body, errors)
		return NewPair[[]any, int64](nil, math.MaxInt64)
	}
	ff, _ := jsons.JsonStringToMap[string, any, JsonObject](body)
	list := ff["statuses"].([]any)
	maxId := cast_x.ToInt64(ff["max_id"])

	return NewPair(list, maxId)
}

func getCookie() string {
	return cookie.GetCookieStrByHost(".weibo.com", "weibo.com")
}
func filter(args map[string]any, groupId int64, word string, day time.Time, uids ...int64) (bson.D, *options.FindOptions) {
	size, ok := args["size"]
	if !ok {
		size = 100
	}

	filter := bson.D{bson.E{Key: "invalid", Value: 0}}
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(cast_x.ToInt64(size))
	if groupId > 0 {
		filter = append(filter, bson.E{Key: "groupId", Value: groupId})
	}
	if len(word) > 0 {
		filter = append(filter, bson.E{Key: "text", Value: "转发微博"})
		return filter, opts
	}

	if day.Year() > 2000 {
		day = day.Add(time.Hour * -8)
		filter = append(filter, bson.E{Key: "created_at", Value: bson.M{
			"$gte": day, "$lt": day.AddDate(0, 0, 1),
		}})
		opts = opts.SetSort(bson.M{"created_at": 1}).SetLimit(cast_x.ToInt64(size))
		DeleteLock()
	} else {
		_kv, _, _ := consts.SfGroup.Do("dao_kv_rss_weibo_config", func() (any, error) {
			key := "rss_weibo_config"
			m := commonkvservice.Get(common_kv_dao.BEAN, key)
			return m[key], nil
		})
		kv := _kv.(*model.CommonKeyValue)

		config := rssWeiboConfig{}
		_ = jsons.JsonStringToObject(kv.CommonValue, &config)
		ignore := config.Ignore

		if len(uids) == 1 && uids[0] == 0 {
			if len(ignore) > 0 {
				filter = append(filter, bson.E{Key: "user.id", Value: bson.M{"$nin": ignore}})
			}
		} else {
			filter = append(filter, bson.E{Key: "user.id", Value: bson.M{"$in": uids}})
			for _, uid := range uids {
				if !slices.Contains(ignore, uid) {
					common_kv_dao.JsonArrayAppend(kv.ID, "$.ignore", uid)
				}
			}
		}
	}
	return filter, opts
}
func Rss(args map[string]any, groupId int64, word string, day time.Time, uids ...int64) string {
	filter, opts := filter(args, groupId, word, day, uids...)

	if cast_x.ToBool(args["deleteBefore"]) {
		DeleteLock()
	}
	if cast_x.ToBool(args["pullBefore"]) {
		PullByGroupLock()
	}

	cursor, _ := collection.Find(context.Background(), filter, opts)
	defer cursor.Close(context.Background())

	s0 := ""
	for cursor.Next(context.Background()) {
		var jo JsonObject
		_ = cursor.Decode(&jo)
		s0 += a(jo)
	}
	s := `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
        <title><![CDATA[网络傻事]]></title>
%s
    </channel>
</rss>
`
	return fmt.Sprintf(s, s0)
}

func a(jo JsonObject) string {
	_id := cast_x.ToInt64(jo["_id"])
	title := jo["text_raw"].(string)
	_contentHtml := contentHtml(jo)
	retweetId := jo["retweeted_status"]
	if retweetId != nil {
		var retweet JsonObject
		_ = collection.FindOne(context.TODO(), bson.M{"_id": cast_x.ToInt64(retweetId)}).Decode(&retweet)
		if retweet != nil {
			i := retweet["user"]
			uName := ""
			if i != nil {
				uName = i.(JsonObject)["name"].(string)
			}
			t := uName + strings.Repeat(consts.Ensp, 4) + time.UnixMilli(cast_x.ToInt64(retweet["created_at"])).Format(time.DateTime) +
				common_x.IfThen(cast_x.ToBool(retweet["read"]), "<br/>pass", "")
			_contentHtml = fmt.Sprintf("%s<hr/>转发自:  @%s<br/>%s", _contentHtml, t, contentHtml(retweet))
		}
	}
	deleteUrl := addDelete(_id)
	_contentHtml = fmt.Sprintf("%s<br/><br/>%s<br/><br/>%s", deleteUrl, _contentHtml, deleteUrl)
	createdAt := time.UnixMilli(cast_x.ToInt64(jo["created_at"])).Format(time.RFC3339)
	user := jo["user"].(JsonObject)
	userId := cast_x.ToInt64(user["id"])
	screenName := user["name"]
	url := fmt.Sprintf("https://weibo.com/%v/%v", userId, jo["mblogid"])
	return fmt.Sprintf(
		`
           <item>
               <title>
                   <![CDATA[ %v ]]>
               </title>
               <description>
                   <![CDATA[ %v ]]>
               </description>
               <pubDate>%v</pubDate>
               <link>%v</link>
               <guid>%v</guid>
               <author>
                   <![CDATA[ %v ]]>
               </author>
           </item>
`, title[0:maths_x.Min(10, len(title))], _contentHtml, createdAt, url, _id, screenName)
}

func addDelete(_id int64) string {
	format := `<a href="http://` + consts.ServiceHost + `:58090/weibo/rss/delete/%v">删除<a/>`
	return fmt.Sprintf(format, _id)
}

func contentHtml(jo JsonObject) string {
	text := jo["text"].(string)
	text = strings.ReplaceAll(text, "//<a ", "<br/>//<a ")
	text = strings.ReplaceAll(text, "//@", "<br/>//@")
	text = strings.ReplaceAll(text, "\n", "<br/>")
	text = strings.ReplaceAll(text, "。”", "。”<br/>")

	text = aa(text)

	picList := ""
	picUrls := jo["pic_ids"].(bson.A)
	for i, url := range picUrls {
		pUrl, _ := nets_x.UrlAddQuery("http://"+consts.ImgHost+":58090/redirect", map[string]any{
			"url": url.(string),
		})
		picList += "<br/>" + strconv.Itoa(i) + "<br/>"
		picList += "<img vspace=\"8\" hspace=\"4\" style=\"\" src=\"" + pUrl.String() + "\" referrerpolicy=\"no-referrer\">"
	}
	return text + picList
}

var sampleRegexp = regexp.MustCompile(`<a\s+[^>]*href="(.*?)".*?>(.*?)<\/a>`)

func aa(text string) string {
	allString := sampleRegexp.FindAllStringSubmatch(text, -1)
	format := sampleRegexp.ReplaceAllString(text, "%s")
	r := make([]any, 0)
	for _, ss := range allString {
		if !(len(ss) > 2) || !strings.Contains(ss[2], "查看图片") {
			r = append(r, ss[0])
			continue
		}
		pUrl, _ := nets_x.UrlAddQuery("http://"+consts.ImgHost+":58090/redirect", map[string]any{
			"url": ss[1],
		})
		a := "<img vspace=\"8\" hspace=\"4\" style=\"\" src=\"" + pUrl.String() + "\" referrerpolicy=\"no-referrer\">"
		r = append(r, a)
	}
	return fmt.Sprintf(format, r...)
}

func DeleteById(id int64) int64 {
	r, _ := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return r.DeletedCount
}
