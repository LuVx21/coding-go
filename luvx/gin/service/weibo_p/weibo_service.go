package weibo_p

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/config"
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/dao/freshrss_dao"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/model"
	"luvx/gin/service"
	commonkvservice "luvx/gin/service/common_kv"
	"luvx/gin/service/cookie"
	"luvx/gin/service/rpc"
	"luvx/gin/service/rss"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/icloudza/fxjson"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/common_x/t"
	"github.com/luvx21/coding-go/coding-common/iterators"
	"github.com/luvx21/coding-go/coding-common/jsons"
	"github.com/luvx21/coding-go/coding-common/maps_x"
	"github.com/luvx21/coding-go/coding-common/nets_x"
	"github.com/luvx21/coding-go/coding-common/sets"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	"github.com/luvx21/coding-go/luvx_service_sdk/proto_gen/proto_kv"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	COMMON_KEY = "rss_weibo_config"
)

var (
	fields = []string{"_id", "user", "user_id", "mblogid", "created_at", "text", "retweeted_status", "pic_ids",
		"invalid", "read", "extra", "groupId", "category",
		// "mix_media_info", "text_raw", "page_info",
	}
	collection    = mongo_dao.WeiboFeedCol
	ignoreRssUids = []int64{}
)

var (
	sourceGroup    = strings.Split(config.Viper.GetString("rss.weibo.sourceGroup"), ",")
	saveImageGroup = strings.Split(config.Viper.GetString("rss.weibo.saveImageGroup"), ",")

	sampleRegexp = regexp.MustCompile(`<a\s+[^>]*href="(.*?)".*?>(.*?)<\/a>`)
	replacer     = strings.NewReplacer(
		"\n", "<br/>",
		"。”", "。”<br/>",
		"//@", "<br/>//@",
		"//<a ", "<br/>//<a ",
	)
)

func PullHotBand() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		// log.Infoln("请求:", r.URL.String())
		r.Headers.Add("Referer", "https://weibo.com/hot/search")
	})

	c.OnResponse(func(r *colly.Response) {
		ff := make(map[string]any)
		_ = sonic.Unmarshal(r.Body, &ff)

		client := mongo_dao.WeiboHotCol
		bandList := ff["data"].(map[string]any)["band_list"].([]any)
		now := times_x.TimeNowDate()
		toInsert := make([]any, 0)
		words := slices_x.Transfer(func(item any) any { return item.(map[string]any)["word"] }, bandList...)
		rows, _ := mongodb.RowsMap(context.TODO(), client, bson.M{"word": bson.M{"$in": words}})
		rowsMap := lo.KeyBy(*rows, func(m bson.M) string { return cast_x.ToString(m["word"]) })

		for _, v := range bandList {
			vv := v.(map[string]any)
			rank := vv["realpos"]
			if rank == nil {
				continue
			}
			word := vv["word"]
			result := rowsMap[cast_x.ToString(word)]
			if result != nil {
				rankMap := mongodb.DM(result["rankMap"].(bson.D))
				oldRank := maps_x.GetOrDefault(rankMap, now, "99")
				if cast_x.ToInt(oldRank) > cast_x.ToInt(rank) {
					rankMap[now] = cast_x.ToString(rank)
					update := bson.M{"$set": bson.M{
						"rankMap":  rankMap,
						"category": vv["category"],
					}}
					_, _ = client.UpdateOne(context.TODO(), bson.M{"word": word}, update)
				}
			} else {
				d := bson.M{
					// "_class":   "org.luvx.boot.tools.dao.mongo.weibo.HotBand",
					"_id":      consts.IdWorker.NextId(),
					"word":     word,
					"category": vv["category"],
					"rankMap":  map[string]string{now: cast_x.ToString(rank)},
				}
				toInsert = append(toInsert, d)
			}
			// fmt.Println(i+1, word, result == nil)
		}
		_, _ = client.InsertMany(context.TODO(), toInsert)
		// fmt.Println("Visited", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Errorln("请求异常:", r.Request.URL.RequestURI(), err.Error())
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
	_ = collection.FindOne(context.TODO(), bson.M{"user_id": uid}, opts).Decode(&latest)

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
	collection.InsertMany(context.TODO(), arr, options.InsertMany().SetOrdered(false))
}

func PullByGroupLock() {
	service.RunnerLocker().LockRun("拉取分组微博", time.Minute*3, func() {
		for _, groupId := range sourceGroup {
			PullByGroup(cast_x.ToInt64(groupId))
		}
	})
}

func PullByGroup(groupId int64) {
	opts := options.FindOne().SetSort(bson.M{"_id": -1})
	var latest bson.M
	_ = collection.FindOne(context.TODO(), bson.M{"groupId": groupId}, opts).Decode(&latest)

	var cursor int64 = 0
	iterator := iterators.NewCursorIterator(
		cursor, false,
		func(_cursor int64) t.Pair[[]any, int64] {
			return requestPageOfGroup(groupId, _cursor)
		},
		func(curId int64, p t.Pair[[]any, int64]) int64 {
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
		func(p t.Pair[[]any, int64]) []any {
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
			f["groupId"] = groupId
			feeds = append(feeds, f)
		}
		parseAndSaveFeed(feed, false)
		feed["groupId"] = groupId
		feeds = append(feeds, feed)
	})

	runs.Go(func() {
		if rpc.KvRpcClient == nil || !slices.Contains(saveImageGroup, strconv.FormatInt(groupId, 10)) {
			return
		}
		urls := slices_x.FlatMap(feeds, func(feed any) []string { return feed.(map[string]any)["pic_ids"].([]string) })
		toSave := sets.NewSet(urls...)
		for _url := range *toSave {
			_, err := (*rpc.KvRpcClient).Get(context.TODO(), &proto_kv.Key{Key: _url})
			if err == nil {
				continue
			}

			_ = consts.RateLimiter.Wait(context.TODO())
			_, body, errors := consts.GoRequest.Get(_url).EndBytes()
			if len(errors) > 0 {
				continue
			}
			_, err = (*rpc.KvRpcClient).Put(context.TODO(), &proto_kv.PutRequest{Entry: &proto_kv.Entry{Key: _url, Value: body}, Expire: int64(7 * 24 * time.Hour.Seconds())})
			if err != nil {
				slog.Error("rpc调用错误", "err", err.Error(), "rpc client", rpc.KvRpcClient == nil)
			}
		}
	})

	_, _ = collection.InsertMany(context.TODO(), feeds, options.InsertMany().SetOrdered(false))
}

func parseAndSaveFeed(feed map[string]any, retweeted bool) int64 {
	id := cast_x.ToInt64(feed["id"])
	feed["_id"] = id
	if cast_x.ToBool(feed["isLongText"]) {
		feed["text"] = requestLongText(feed["mblogid"].(string))
	}
	category := extractTagsManual(feed["text"].(string))
	if len(category) > 0 {
		feed["category"] = category
	}

	picUrl := make([]string, 0)
	i2 := feed["pic_ids"]
	if i2 != nil {
		if picIds, b := slices_x.IsNotEmpty(i2.([]any)); b {
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
		feed["user_id"] = cast_x.ToInt64(user["id"])
		feed["user"] = map[string]any{"name": user["screen_name"]}
	} else {
		feed["user_id"] = 0
		feed["user"] = map[string]any{"name": ""}
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
	_, body, _ := requestWeibo("https://weibo.com/ajax/statuses/longtext", m, map[string]string{"Cookie": getCookie()})

	return gjson.Get(body, "data.longTextContent").String()
}

func requestPageOfUser(uid int64, cursor int) []any {
	m := map[string]any{
		"uid":     uid,
		"page":    cursor,
		"feature": 0,
	}
	_, body, _ := requestWeibo("https://weibo.com/ajax/statuses/mymblog", m, map[string]string{"Cookie": getCookie()})

	ff, _ := jsons.JsonStringToMap[string, any, map[string]any](body)
	i := ff["data"]
	if i == nil {
		return make([]any, 0)
	}
	list := i.(map[string]any)["list"].([]any)
	return list
}

func requestPageOfGroup(groupId int64, cursor int64) t.Pair[[]any, int64] {
	m := map[string]any{
		"list_id":      groupId,
		"max_id":       cursor,
		"refresh":      4,
		"fast_refresh": 1,
		"count":        25,
	}
	var list []any
	var maxId int64 = 0
	for list == nil {
		r, body, errors := requestWeibo("https://weibo.com/ajax/feed/groupstimeline", m, map[string]string{"Cookie": getCookie()})
		if len(errors) > 0 {
			return t.NewPair[[]any, int64](nil, math.MaxInt64)
		}
		ff, _ := jsons.JsonStringToMap[string, any, bson.M](body)
		if ff["statuses"] == nil {
			slog.Error("异常响应数据", "url", r.Request.URL.String(), "body", body)
		} else {
			list, maxId = ff["statuses"].([]any), cast_x.ToInt64(ff["max_id"])
		}
	}
	return t.NewPair(list, maxId)
}

func getCookie() string {
	m := cookie.GetCookieStrByHost(".weibo.com", "weibo.com")
	return m[".weibo.com"] + "; " + m["weibo.com"]
}
func filter(args map[string]any, groupId int64, word string, day time.Time, uids ...int64) (bson.M, options.Lister[options.FindOptions]) {
	size, ok := args["size"]
	if !ok {
		size = 100
	}

	filter := bson.M{"invalid": 0}
	opts := options.Find().SetSort(bson.M{"created_at": -1}).SetLimit(cast_x.ToInt64(size))
	if cast_x.ToBool(args["asc"]) {
		opts = opts.SetSort(bson.M{"created_at": 1})
	}

	if groupId > 0 {
		filter["groupId"] = groupId
	}
	if len(word) > 0 {
		w := cast_x.ToString(commonkvservice.GetMapFieldValue("common_map", word))
		if len(w) > 0 {
			filter["text"] = bson.M{"$regex": w, "$options": "i"}
			return filter, opts
		}
	}

	if day.Year() > 2000 {
		day = day.Add(time.Hour * -8)
		filter["created_at"] = bson.M{"$gte": day, "$lt": day.AddDate(0, 0, 1)}
		DeleteLock()
	} else {
		if len(ignoreRssUids) == 0 {
			_kv, _, _ := consts.SfGroup.Do("dao_kv_rss_weibo_config", func() (any, error) {
				m := commonkvservice.Get(common_kv_dao.BEAN, COMMON_KEY)
				return m[COMMON_KEY], nil
			})
			kv := _kv.(*model.CommonKeyValue)

			config := rssWeiboConfig{}
			_ = jsons.JsonStringToObject(kv.CommonValue, &config)
			ignoreRssUids = config.Ignore
		}
		if len(uids) == 1 && uids[0] == 0 {
			if len(ignoreRssUids) > 0 {
				filter["user_id"] = bson.M{"$nin": ignoreRssUids}
			}
		} else {
			filter["user_id"] = common_x.IfThen[any](len(uids) == 1, uids[0], bson.M{"$in": uids})
			flag := false
			for _, uid := range uids {
				if !slices.Contains(ignoreRssUids, uid) {
					common_kv_dao.JsonArrayAppend(common_kv_dao.BEAN, COMMON_KEY, "$.ignore", uid)
					flag = true
				}
			}
			if flag {
				ignoreRssUids = ignoreRssUids[:0]
			}
		}
	}
	return filter, opts
}
func Rss(c *gin.Context, args map[string]any, groupId int64, word string, day time.Time, uids ...int64) string {
	k := strconv.FormatInt(uids[0], 10)
	filter, opts := filter(args, groupId, word, day, uids...)

	if cast_x.ToBool(args["deleteBefore"]) {
		DeleteLock()
	}
	if cast_x.ToBool(args["pullBefore"]) {
		PullByGroupLock()
	}

	p := common_x.RunWithTimeReturn(k+":weibo_rss_1", func() t.Pair[*[]bson.M, error] {
		cursor, err := mongodb.RowsMap(context.TODO(), collection, filter, opts)
		return t.NewPair(cursor, err)
	})
	cursor, err := p.Unpack()
	if err != nil || len(*cursor) == 0 {
		return rss.ToRssXml(nil, "网络傻事")
	}
	ids, retweetedIds := make([]string, 0), make([]int64, 0)
	for i := range *cursor {
		jo := (*cursor)[i]
		ids = append(ids, cast_x.ToString(jo["_id"]))
		if jo["retweeted_status"] != nil {
			retweetedIds = append(retweetedIds, cast_x.ToInt64(jo["retweeted_status"]))
		}
	}
	var rowsMap = map[int64]bson.M{}
	if len(retweetedIds) > 0 {
		p := common_x.RunWithTimeReturn(k+":weibo_rss_2", func() t.Pair[*[]bson.M, error] {
			rows, err := mongodb.RowsMap(context.TODO(), collection, bson.M{"_id": bson.M{"$in": retweetedIds}})
			return t.NewPair(rows, err)
		})
		rows, _ := p.Unpack()
		rowsMap = lo.KeyBy(*rows, func(m bson.M) int64 { return cast_x.ToInt64(m["_id"]) })
	}

	aa := common_x.IfThenGet(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100) < 60,
		func() []string { return freshrss_dao.ExistedGuids("%"+c.Request.URL.Path+"%", ids) },
		func() []string { return []string{} },
	)
	existedGuids := sets.NewSet(aa...)

	feeds := make([]*rss.RssItem, 0)
	for i := range *cursor {
		jo := (*cursor)[i]
		if existedGuids.Contains(cast_x.ToString(jo["_id"])) {
			continue
		}
		feeds = append(feeds, rssItem(jo, rowsMap[cast_x.ToInt64(jo["retweeted_status"])]))
	}

	return rss.ToRssXml(feeds, "网络傻事")
}

func rssItem(jo, retweet bson.M) *rss.RssItem {
	_id := cast_x.ToInt64(jo["_id"])
	// title := jo["text_raw"].(string)
	var categories bson.A
	if cs, ok := jo["category"]; ok && cs != nil {
		categories = append(categories, cs.(bson.A)...)
	}
	_contentHtml := contentHtml(jo)
	if retweet != nil {
		if cs, ok := retweet["category"]; ok && cs != nil {
			categories = append(categories, cs.(bson.A)...)
		}

		user, retweetUrl, uName := retweet["user"], "", ""
		if user != nil {
			umap := mongodb.DM(user.(bson.D))
			uId := cast_x.ToString(retweet["user_id"])
			retweetUrl = fmt.Sprintf("<a href=\"https://weibo.com/%s/%s\">转发自</a>", uId, retweet["mblogid"])
			uName = fmt.Sprintf("<a href=\"https://weibo.com/u/%s\">@%s</a>", uId, umap["name"].(string))
		}
		t := uName + strings.Repeat(consts.Ensp, 4) + time.UnixMilli(cast_x.ToInt64(retweet["created_at"])).Format(time.DateTime) +
			common_x.IfThen(cast_x.ToBool(retweet["read"]), "<br/>pass", "")
		_contentHtml = fmt.Sprintf("%s<hr/>%s:  %s<br/>%s", _contentHtml, retweetUrl, t, contentHtml(retweet))
	}
	deleteUrl := addDelete(_id)
	_contentHtml = fmt.Sprintf("%s<br/><br/>%s", _contentHtml, deleteUrl)
	createdAt := time.UnixMilli(cast_x.ToInt64(jo["created_at"])).Format(time.RFC3339)
	screenName := mongodb.DM(jo["user"].(bson.D))["name"]
	url := fmt.Sprintf("https://weibo.com/%v/%v", cast_x.ToInt64(jo["user_id"]), jo["mblogid"])
	rssItem := rss.RssItem{
		Title:       "title",
		Description: _contentHtml,
		PubDate:     createdAt,
		Link:        url,
		Guid:        cast_x.ToString(_id),
		Author:      screenName.(string),
		Categories:  slices_x.Transfer(func(a any) string { return cast_x.ToString(a) }, categories...),
	}
	return &rssItem
}

func addDelete(_id int64) string {
	format := `<a href="http://` + consts.ServiceDomain + `/rss/delete/%s/%v?real=true">删除<a/>`
	return fmt.Sprintf(format, mongo_dao.COL_NAME_weibo_feed, _id)
}

func contentHtml(jo bson.M) string {
	text := jo["text"].(string)
	text = replacer.Replace(text)

	text = aa(text)

	picUrls := jo["pic_ids"].(bson.A)
	picList, pc := "", strconv.Itoa(len(picUrls))
	for i, _url := range picUrls {
		url := roundImgCdn(_url.(string))
		picList += "<br/>" + strconv.Itoa(i+1) + "/" + pc + "<br/>"
		picList += "<img style=\"margin: 8px 4px;width:300px\" src=\"" + url + "\" referrerpolicy=\"no-referrer\">"
	}
	return text + picList
}

func aa(text string) string {
	allString := sampleRegexp.FindAllStringSubmatch(text, -1)
	format := sampleRegexp.ReplaceAllString(text, "%s")
	r := make([]any, 0)
	for _, ss := range allString {
		if !(len(ss) > 2) || !strings.Contains(ss[2], "查看图片") {
			r = append(r, ss[0])
			continue
		}
		url := roundImgCdn(ss[1])
		a := "<img style=\"margin: 8px 4px;width:300px\" src=\"" + url + "\" referrerpolicy=\"no-referrer\">"
		r = append(r, a)
	}
	return fmt.Sprintf(format, r...)
}

// 各种图片CDN
func roundImgCdn(_url string) string {
	if !config.GetSwitch("weibo.imgCdn") {
		return _url
	}
	i := rand.Intn(4)
	if i > 0 {
		u2, _ := url.Parse(_url)
		switch i {
		case 1:
			return "https://cdn.ipfsscan.io/weibo" + u2.Path
		case 2:
			return "https://cdn.cdnjson.com/" + u2.Host + u2.Path
		case 3:
			return "https://i" + strconv.Itoa(rand.Intn(4)) + ".wp.com/" + u2.Host + u2.Path
		}
	}

	u1, _ := nets_x.UrlAddQuery(consts.ImgRedirectUrlPrefix, map[string]any{"url": _url})
	return u1.String()
}

func getGroupMember(groupId int64) []int64 {
	data := func(c int) []int64 {
		m := map[string]any{
			"list_id": groupId,
			"page":    c,
		}
		_, body, _ := requestWeibo("https://weibo.com/ajax/profile/getGroupMembers", m, map[string]string{"Cookie": getCookie()})
		node := fxjson.FromString(body)
		nextCursor := node.Get("data.next_cursor").IntOr(0)
		r := make([]int64, 0)
		if nextCursor == 0 {
			return r
		}
		node.GetPath("data.users").ArrayForEach(func(index int, item fxjson.Node) bool {
			id := item.Get("id").IntOr(0)
			if id > 0 {
				r = append(r, id)
			}
			return true
		})
		return r
	}

	r := make([]int64, 0)
	for i := range 99 {
		ids := data(i + 1)
		if len(ids) == 0 {
			break
		}
		r = append(r, ids...)
	}

	return r
}

func extractTagsManual(s string) []string {
	var tags []string
	var inTag bool
	var tagStart int

	for i := 0; i < len(s); i++ {
		if s[i] != '#' {
			continue
		}
		if inTag {
			// 结束标签
			if i > tagStart+1 { // 避免 ##
				tags = append(tags, s[tagStart+1:i])
			}
			inTag = false
		} else {
			// 开始标签
			tagStart = i
			inTag = true
		}
	}
	return tags
}
