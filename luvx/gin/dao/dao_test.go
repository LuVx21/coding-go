package dao

import (
	"fmt"
	"luvx/gin/dao/common_kv_dao"
	"testing"
)

func Test_01(t *testing.T) {
	m := common_kv_dao.Get(common_kv_dao.INDEX_SPIDER, "foo", "bar")
	for _, e := range m {
		fmt.Println(e.CommonKey)
	}
}
