package rss

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"luvx/gin/common/consts"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/service/common_kv"
	"luvx/gin/service/soup"

	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x/a"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/nosql/mongodb"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Rss(spiderKey string) string {
	filter := bson.M{"spiderKey": spiderKey, "invalid": 0}
	opts := options.Find().SetSort(bson.D{{Key: "pubDate", Value: -1}, {Key: "_id", Value: -1}}).SetLimit(100)
	ms, _ := mongodb.RowsMap(context.Background(), mongo_dao.RssFeedCol, filter, opts)

	result := make([]*RssItem, 0)
	for i := range *ms {
		result = append(result, parse2RssItem((*ms)[i]))
	}
	return ToRssXml(result, spiderKey)
}

func parse2RssItem(m a.JsonObject) *RssItem {
	_id := m["_id"].(int64)
	contents := m["content"].(bson.A)

	contentHtml := ""
	for _, c := range contents {
		url := c.(string)
		if strings.HasPrefix(url, "http") {
			contentHtml += fmt.Sprintf("<img vspace=\"8\" hspace=\"4\" style=\"\" src=\"%s\" referrerpolicy=\"no-referrer\">", url)
		} else {
			contentHtml += "<div>" + url + "</div>"
		}
	}

	deleteUrl := fmt.Sprintf(`<a href="http://`+consts.ServiceDomain+`/rss/delete/%s/%v">删除<a/>`, mongo_dao.COL_NAME_rss_feed, _id)
	contentHtml = deleteUrl + `<br/>` + contentHtml + `<br/>` + deleteUrl
	cs := slices_x.Transfer(func(a any) string { return cast_x.ToString(a) }, (m["categorySet"].(bson.A))...)

	return &RssItem{
		Title:       m["title"].(string),
		Description: contentHtml,
		PubDate:     m["pubDate"].(string),
		Link:        m["url"].(string),
		Guid:        strconv.FormatInt(_id, 10),
		Author:      "未知",
		Categories:  cs,
	}
}

func PullByKey() {
	m := common_kv.Get(8)
	for k, v := range m {
		runs.Go(func() {
			log.Infoln("spider拉取:", k)
			items := spiderIndexPage(k, v.CommonValue)
			_, _ = mongo_dao.RssFeedCol.InsertMany(context.TODO(), slices_x.ToAnySliceE(items...))
		})
	}
}

func spiderIndexPage(key, paramJson string) []soup.PageContent {
	filter := bson.M{
		"spiderKey": key,
		"content":   bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
	}
	opts := options.Find().
		SetProjection(bson.M{
			"url": 1,
			// "content": 1,
		}).
		SetSort(bson.M{"_id": -1}).SetLimit(2000)
	cursor, _ := mongo_dao.RssFeedCol.Find(context.Background(), filter, opts)
	defer cursor.Close(context.Background())

	ignoreIndexItemUrlSet := make([]string, 0, 2000)
	for cursor.Next(context.Background()) {
		m := make(map[string]any)
		_ = cursor.Decode(&m)
		ignoreIndexItemUrlSet = append(ignoreIndexItemUrlSet, m["url"].(string))
	}

	return soup.Of(paramJson).
		SetIgnoreIndexItemUrlSet(ignoreIndexItemUrlSet).
		Visit()
}
