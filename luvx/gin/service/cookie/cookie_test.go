package cookie

import (
	"testing"
)

func Test_01(t *testing.T) {
	cookie := GetCookieByHost(".weibo.com", "weibo.com")
	t.Log(cookie)
}
