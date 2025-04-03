package service

import (
	"fmt"
	"luvx/gin/dao/common_kv_dao"
	"luvx/gin/service/common_kv"
	"luvx/gin/service/cookie"
	"testing"
)

func Test_00(t *testing.T) {
}

func Test_01(t *testing.T) {
	cookie := cookie.GetCookieByHost(".weibo.com", "weibo.com")
	t.Log(cookie)
}

func Test_02(t *testing.T) {
	m := common_kv.Get(common_kv_dao.INDEX_SPIDER, "foo", "bar")
	for k, v := range m {
		fmt.Println(k, v)
	}
}
