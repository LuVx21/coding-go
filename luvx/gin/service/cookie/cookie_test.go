package cookie

import (
	"testing"
)

func Test_01(t *testing.T) {
	cookie := GetCookieFromDb(".weibo.com", "weibo.com")
	t.Log(cookie)
}
