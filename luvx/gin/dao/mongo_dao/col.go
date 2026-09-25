package mongo_dao

import db_mongo "luvx/gin/db/mongo"

const (
	COL_NAME_rss_feed   = "rss_feed"
	COL_NAME_bili_video = "bili_video"
	COL_NAME_weibo_feed = "weibo_feed"
)

var (
	ConfigCol    = db_mongo.GetMainCollection("config")
	CookieCol    = db_mongo.GetMainCollection("cookie")
	UserCol      = db_mongo.GetMainCollection("user")
	WeiboFeedCol = db_mongo.GetMainCollection(COL_NAME_weibo_feed)
	WeiboHotCol  = db_mongo.GetMainCollection("weibo_hot_band")

	BiliVideoCol = db_mongo.GetSlaveCollection(COL_NAME_bili_video)
	RssFeedCol   = db_mongo.GetSlaveCollection(COL_NAME_rss_feed)
)
