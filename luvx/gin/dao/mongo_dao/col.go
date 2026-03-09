package mongo_dao

import "luvx/gin/db"

var (
	ConfigCol    = db.GetMainCollection("config")
	CookieCol    = db.GetMainCollection("cookie")
	UserCol      = db.GetMainCollection("user")
	WeiboFeedCol = db.GetMainCollection("weibo_feed")
	WeiboHotCol  = db.GetMainCollection("weibo_hot_band")

	BiliVideoCol = db.GetSlaveCollection("bili_video")
	RssFeedCol   = db.GetSlaveCollection("rss_feed")
)
