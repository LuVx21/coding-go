package mongo_dao

import "luvx/gin/db"

const (
	COL_NAME_rss_feed   = "rss_feed"
	COL_NAME_bili_video = "bili_video"
	COL_NAME_weibo_feed = "weibo_feed"
)

var (
	ConfigCol    = db.GetMainCollection("config")
	CookieCol    = db.GetMainCollection("cookie")
	UserCol      = db.GetMainCollection("user")
	WeiboFeedCol = db.GetMainCollection(COL_NAME_weibo_feed)
	WeiboHotCol  = db.GetMainCollection("weibo_hot_band")

	BiliVideoCol = db.GetSlaveCollection(COL_NAME_bili_video)
	RssFeedCol   = db.GetSlaveCollection(COL_NAME_rss_feed)
)
