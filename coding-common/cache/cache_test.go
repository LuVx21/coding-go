package cache

import (
	"fmt"
	"maps"
	"testing"

	"github.com/luvx21/coding-go/coding-common/maps_x"
)

type (
	row = map[int]string
)

var (
	source1 = &dataSource[int, string]{name: "数据源1", data: row{1: "a", 2: "b"}}
	source2 = &dataSource[int, string]{name: "数据源2", data: row{3: "c", 4: "d"}}
	source3 = &dataSource[int, string]{name: "数据源3", data: row{5: "e", 6: "f"}}
)

type dataSource[K comparable, V any] struct {
	name string
	data map[K]V
}

func (m *dataSource[K, V]) Get(ids []K) map[K]V {
	fmt.Println("从"+m.name+"获取数据", ids)
	return maps_x.GetByKeys(m.data, ids...)
}
func (m *dataSource[K, V]) Set(data map[K]V) {
	fmt.Println(m.name+"写入前数据", m.data, "待写入数据", data)
	if m.data == nil {
		m.data = make(map[K]V)
	}
	maps.Copy(m.data, data)
	fmt.Println(m.name+"写入后数据", m.data)
}
func (m *dataSource[K, V]) Data() map[K]V {
	return m.data
}

func Test_cache_00(t *testing.T) {
	r := GetFrom([]int{1, 5}, source1, source2, source3)
	fmt.Println(r)
	fmt.Println("------------------------------")
	r = GetFrom([]int{1, 5}, source1, source2, source3)
	fmt.Println(r)
}

func Test_cache_01(t *testing.T) {
	r := GetNonSet([]int{1, 3, 5},
		func(k []int) row { return source1.Data() },
		func(k []int) row { return source2.Data() },
		func(k []int) row { return source3.Data() },
	)
	fmt.Println(r)
}
