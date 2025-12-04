package common_kv

import (
	"fmt"
	"luvx/gin/dao/common_kv_dao"
	"testing"
)

func Test_02(t *testing.T) {
	m := Get(common_kv_dao.INDEX_SPIDER, "foo", "bar")
	for k, v := range m {
		fmt.Println(k, v)
	}
}
