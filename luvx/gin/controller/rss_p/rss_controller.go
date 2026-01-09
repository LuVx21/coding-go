package rss_p

import (
	"context"
	"luvx/gin/common/responsex"
	"luvx/gin/dao/mongo_dao"
	"luvx/gin/db"
	"luvx/gin/service/rss"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"go.mongodb.org/mongo-driver/bson"
)

func Rss(c *gin.Context) {
	rss := rss.Rss(c.Param("spiderKey"))
	c.Header("Content-Type", "application/xml;charset=UTF-8")
	c.String(http.StatusOK, rss)
}

func PullByKey(c *gin.Context) {
	rss.PullByKey()
	c.String(http.StatusOK, "ok")
}

func DeleteById(c *gin.Context) {
	source, id := c.Param("source"), cast_x.ToInt64(c.Param("id"))
	realDel := c.Query("real")
	if source == "" {
		responsex.R(c, "不存在的source")
	}

	cli := db.GetCollectionByName(source)
	if cast_x.ToBool(realDel) {
		n := mongo_dao.DeleteById(cli, id)
		responsex.R(c, map[string]any{"delete": n})
	} else {
		update := bson.M{"$set": bson.M{
			"invalid": 1,
		}}
		one, _ := cli.UpdateOne(context.TODO(), bson.M{"_id": id, "invalid": 0}, update)
		responsex.R(c, one)
	}
}
