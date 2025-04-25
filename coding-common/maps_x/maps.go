package maps_x

import (
	"fmt"
	"reflect"
	"strings"

	maps0 "maps"

	"github.com/luvx21/coding-go/coding-common/reflects"
	"golang.org/x/exp/maps"
)

func Merge[M ~map[K]V, K comparable, V any](m1, m2 M, replace bool) M {
	r := make(M, len(m1)+len(m2))
	maps0.Copy(r, m1)
	for sk, sv := range m2 {
		if _, ok := r[sk]; !ok || replace {
			r[sk] = sv
		}
	}
	return r
}

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return Filter(m, func(k K, v V) bool { return true })
}

func Filter[M ~map[K]V, K comparable, V any](m M, f func(K, V) bool) M {
	r := make(M, len(m))
	for k, v := range m {
		if f(k, v) {
			r[k] = v
		}
	}
	return r
}

func Join[M ~map[K]V, K comparable, V any](m M, kvLink, eLink string) string {
	return JoinMapper(m, kvLink, eLink, nil, nil)
}

// JoinMapper kvLink: kv连接符, eLink: entry连接符, keyMapper,valueMapper: k,v的映射器
func JoinMapper[M ~map[K]V, K comparable, V any](m M,
	kvLink, eLink string,
	keyMapper func(K) string, valueMapper func(V) string) string {
	var sb strings.Builder
	var result string
	keys := maps.Keys(m)
	for i, k := range keys {
		v := m[k]
		var kk, vv any = k, v
		if keyMapper != nil {
			kk = keyMapper(k)
		}
		if valueMapper != nil {
			vv = valueMapper(v)
		}
		sb.WriteString(fmt.Sprintf("%v%s%v", kk, kvLink, vv))
		if i < len(keys)-1 {
			result += eLink
			sb.WriteString(eLink)
		}
	}
	return sb.String()
}

func RemoveIf[M ~map[K]V, K comparable, V any](m M, f func(K, V) bool) {
	for k, v := range m {
		if f(k, v) {
			delete(m, k)
		}
	}
}

func GetOrDefault[M ~map[K]V, K comparable, V any](m M, k K, defaultV V) V {
	if v, exist := m[k]; exist {
		return v
	} else {
		return defaultV
	}
}

func Compute[M ~map[K]V, K comparable, V any](m M, k K, f func(K, V) V) V {
	oldValue, exist := m[k]
	newValue := f(k, oldValue)

	if reflects.IsNil(newValue) {
		if !reflects.IsNil(oldValue) || exist {
			delete(m, k)
		}
		var zero V
		return zero
	}
	m[k] = newValue
	return newValue
}

func ComputeIfAbsent[M ~map[K]V, K comparable, V any](m M, k K, f func(K) V) V {
	oldValue, exist := m[k]
	if !exist || reflects.IsNil(oldValue) {
		newValue := f(k)
		if !reflects.IsNil(newValue) {
			m[k] = newValue
			return newValue
		}
	}
	return oldValue
}

func ComputeIfPresent[M ~map[K]V, K comparable, V any](m M, k K, f func(K, V) V) V {
	oldValue, exist := m[k]
	if exist {
		newValue := f(k, oldValue)
		if reflects.IsNil(newValue) {
			delete(m, k)
			var zero V
			return zero
		}
		m[k] = newValue
		return newValue
	}
	return oldValue
}

func GetInt[M ~map[K]any, K comparable](m M, key K, _default int) (int, error) {
	val, err := GetByKey(m, key, reflect.Int, _default)
	if err != nil {
		return _default, err
	}

	return val.(int), err
}

func GetInt64[M ~map[K]any, K comparable](m M, key K, _default int64) (int64, error) {
	val, err := GetByKey[M, K](m, key, reflect.Int64, _default)
	if err != nil {
		return _default, err
	}

	return val.(int64), err
}

func GetFloat[M ~map[K]any, K comparable](m M, key K, _default float64) (float64, error) {
	val, err := GetByKey(m, key, reflect.Float64, _default)
	if err != nil {
		return _default, err
	}

	return val.(float64), err
}

func GetMap[M ~map[K]any, K, NK comparable, NV any](m M, key K, _default map[NK]NV) (map[NK]NV, error) {
	val, err := GetByKey(m, key, reflect.Map, _default)
	if err != nil {
		return _default, err
	}

	return val.(map[NK]NV), nil
}

func GetString[M ~map[K]any, K comparable](m M, key K, _default string) (string, error) {
	val, err := GetByKey(m, key, reflect.String, _default)
	if err != nil {
		return _default, err
	}

	return val.(string), err
}

func GetSlice[M ~map[K]any, K comparable, E any](m M, key K, _default []E) ([]E, error) {
	val, err := GetByKey(m, key, reflect.Slice, _default)
	if err != nil {
		return _default, err
	}
	return val.([]E), err
}

func GetInterface[M ~map[K]any, K comparable](m M, key K, _default any) (any, error) {
	val, ok := m[key]
	if !ok {
		val = _default
	}

	return val, nil
}

func GetByKey[M ~map[K]any, K comparable](m M, key K, vType reflect.Kind, _default any) (any, error) {
	if m == nil || len(m) == 0 {
		return _default, nil
	}

	val, ok := m[key]
	if !ok {
		return _default, nil
	}

	vk := reflect.ValueOf(val).Kind()
	if vk != vType {
		return nil, fmt.Errorf("非法值类型-> 实际:%s 预期:%s", vk.String(), vType.String())
	}

	return val, nil
}
