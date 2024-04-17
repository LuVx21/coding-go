package service

import (
    "fmt"
    "luvx/gin/service/common_kv"
    "testing"
)

func Test_00(t *testing.T) {
}

func Test_01(t *testing.T) {
    cookie := GetCookieByHost(".weibo.com", "weibo.com")
    t.Log(cookie)
}

func Test_02(t *testing.T) {
    m := common_kv.Get(common_kv.INDEX_SPIDER, "foo", "bar")
    for k, v := range m {
        fmt.Println(k, v)
    }
}
