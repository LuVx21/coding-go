package cache

import (
	"maps"

	"github.com/luvx21/coding-go/coding-common/sets"
)

type (
	Get[K comparable, V any] = func([]K) map[K]V
)

// MultiDataRetrievable 多数据源获取数据, K 具有唯一性且有序的id, V 数据
type MultiDataRetrievable[K comparable, V any] interface {
	// 获取数据
	Get(ids []K) map[K]V
	// 写入数据
	Set(data map[K]V)
}

// GetNonSet 不缓存结果
func GetNonSet[K comparable, V any](ids []K, sources ...Get[K, V]) map[K]V {
	if len(ids) == 0 || len(sources) == 0 {
		return nil
	}
	r, leftIds := make(map[K]V, len(ids)), sets.NewSet[K]()

	data := sources[0](ids)
	for j := range ids {
		id := ids[j]
		if v, ok := data[id]; ok {
			r[id] = v
		} else {
			leftIds.Add(id)
		}
	}

	nextData := GetNonSet(leftIds.ToSlice(), sources[1:]...)
	if len(nextData) > 0 {
		maps.Copy(r, nextData)
	}
	return r
}

// GetFrom 递归从多个数据源获取数据, 取得的数据会放入前一数据源
func GetFrom[K comparable, V any](ids []K, sources ...MultiDataRetrievable[K, V]) map[K]V {
	if len(ids) == 0 || len(sources) == 0 {
		return nil
	}
	r, leftIds := make(map[K]V, len(ids)), sets.NewSet[K]()

	source := sources[0]
	data := source.Get(ids)
	for j := range ids {
		id := ids[j]
		if v, ok := data[id]; ok {
			r[id] = v
		} else {
			leftIds.Add(id)
		}
	}

	nextData := GetFrom(leftIds.ToSlice(), sources[1:]...)
	if len(nextData) > 0 {
		source.Set(nextData)
		maps.Copy(r, nextData)
	}
	return r
}
