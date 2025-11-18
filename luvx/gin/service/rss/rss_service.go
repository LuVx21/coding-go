package rss

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"luvx/gin/common/consts"
	"luvx/gin/common/responsex"
	"luvx/gin/db"
	"luvx/gin/service"
	"luvx/gin/service/common_kv"
	"luvx/gin/service/soup"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x/alias_x"
	"github.com/luvx21/coding-go/coding-common/common_x/runs"
	"github.com/luvx21/coding-go/coding-common/slices_x"
	"github.com/luvx21/coding-go/infra/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ()

var (
	collection = db.MongoDatabase.Collection("rss_feed")
)

func Rss(spiderKey string) string {
	filter := bson.M{"spiderKey": spiderKey, "invalid": 0}
	opts := options.Find().SetSort(bson.D{{Key: "pubDate", Value: -1}, {Key: "_id", Value: -1}}).SetLimit(100)
	cursor, _ := collection.Find(context.Background(), filter, opts)
	defer cursor.Close(context.Background())

	result := make([]*RssItem, 0)
	for cursor.Next(context.Background()) {
		var jo alias_x.JsonObject
		_ = cursor.Decode(&jo)
		result = append(result, parse2RssItem(jo))
	}
	return ToRssXml(result, spiderKey)
}

func parse2RssItem(m alias_x.JsonObject) *RssItem {
	_id := m["_id"].(int64)
	contents := m["content"].(primitive.A)

	contentHtml := ""
	for _, c := range contents {
		url := c.(string)
		if strings.HasPrefix(url, "http") {
			contentHtml += fmt.Sprintf("<img vspace=\"8\" hspace=\"4\" style=\"\" src=\"%s\" referrerpolicy=\"no-referrer\">", url)
		} else {
			contentHtml += "<div>" + url + "</div>"
		}
	}

	deleteUrl := fmt.Sprintf(`<a href="http://`+consts.ServiceHost+`:58090/rss/delete/%v">删除<a/>`, _id)
	contentHtml = deleteUrl + `<br/>` + contentHtml + `<br/>` + deleteUrl

	return &RssItem{
		Title:       m["title"].(string),
		Description: contentHtml,
		PubDate:     m["pubDate"].(string),
		Link:        m["url"].(string),
		Guid:        strconv.FormatInt(_id, 10),
		Author:      "未知",
	}
}

func DeleteById(c *gin.Context) {
	id := cast_x.ToInt64(c.Param("id"))
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "invalid", Value: 1},
		}},
	}
	one, _ := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	responsex.R(c, one)
}

func PullByKey() {
	service.RunnerLocker.LockRun("rss_spider", time.Minute*10, func() {
		m := common_kv.Get(8)
		for k, v := range m {
			runs.Go(func() {
				logs.Log.Infoln("spider拉取:", k)
				items := spiderIndexPage(k, v.CommonValue)
				_, _ = collection.InsertMany(context.TODO(), slices_x.ToAnySliceE(items...))
			})
		}
	})
}

func spiderIndexPage(key, paramJson string) []soup.PageContent {
	filter := bson.M{
		"spiderKey": key,
		"content":   bson.M{"$exists": true, "$not": bson.M{"$size": 0}},
	}
	opts := options.Find().
		SetProjection(bson.D{
			{Key: "url", Value: 1},
			// {"content", 1},
		}).
		SetSort(bson.M{"_id": -1}).SetLimit(2000)
	cursor, _ := collection.Find(context.Background(), filter, opts)
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
