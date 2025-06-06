package jsonpicker

import (
	"encoding/json"
	"reflect"
	"slices"
)

type KeyStringMap struct {
	innerMap map[string]any
}

type KSMapEntry struct {
	Key   any
	Value any
}

type JsonValue struct {
	value any
}

func (j JsonValue) IsMap() bool {
	return reflect.TypeOf(j.value).Kind() == reflect.Map
}

func (j JsonValue) IsSliceOrArray() bool {
	return reflect.TypeOf(j.value).Kind() == reflect.Slice || reflect.TypeOf(j.value).Kind() == reflect.Array
}

func (j JsonValue) Get(key string) IJsonValue {
	if j.IsMap() {
		m := ConvKSMap(j.value)
		if vv := m.Get(key); vv != nil {
			return JsonValue{value: vv.Value()}
		} else {
			return nil
		}
	}
	return nil
}

func deepGet(j IJsonValue, key string) IJsonValue {
	if v := j.Get(key); v != nil {
		return v
	} else {
		m := ConvKSMap(j.Value())
		values := m.Values()
		maps := values.Filter(func(v any) bool {
			return reflect.TypeOf(v).Kind() == reflect.Map
		})
		if maps.Len() == 0 {
			return nil
		} else {
			for i := range maps.Len() {
				jj := JsonValue{value: maps.GetAt(i)}
				rr := deepGet(jj, key)
				if rr != nil && rr.Value() != nil {
					return rr
				}
			}
		}
	}
	return nil
}

func (j JsonValue) DeepGet(key string) IJsonValue {
	if j.IsMap() {
		return deepGet(j, key)
	}
	return nil
}

func (j JsonValue) GetAt(index int) IJsonValue {
	if j.IsSliceOrArray() {
		m := ConvSlice(j.value)
		vv := m.GetAt(index)
		return JsonValue{value: vv}
	}
	return nil
}

func (j JsonValue) ContainsKey(key string) bool {
	if j.IsMap() {
		m := ConvKSMap(j.value)
		return m.ContainsKey(key)
	}
	return false
}

func (j JsonValue) Contains(v any) bool {
	if j.IsSliceOrArray() {
		s := ConvSlice(j.value)
		return s.Contains(v)
	}
	return false
}

func (j JsonValue) Value() any {
	return j.value
}

func (m KeyStringMap) Get(key string) IJsonValue {
	if v, ok := m.innerMap[key]; ok {
		return JsonValue{value: v}
	} else {
		return nil
	}
}

func mapDeepGet(m *KeyStringMap, key string) IJsonValue {
	if v := m.Get(key); v != nil {
		return v
	} else {
		values := m.Values()
		maps := values.Filter(func(v any) bool {
			return reflect.TypeOf(v).Kind() == reflect.Map
		})
		if maps.Len() == 0 {
			return nil
		} else {
			for i := range maps.Len() {
				jj := ConvKSMap(maps.GetAt(i))
				rr := mapDeepGet(jj, key)
				if rr != nil && rr.Value() != nil {
					return rr
				}
			}
		}
	}
	return nil
}

func (m *KeyStringMap) DeepGet(key string) IJsonValue {
	if v, ok := m.innerMap[key]; ok {
		return JsonValue{value: v}
	} else {
		return mapDeepGet(m, key)
	}
}

func (m *KeyStringMap) Set(key string, value any) {
	m.innerMap[key] = value
}

func (m KeyStringMap) ContainsKey(key string) bool {
	if _, ok := m.innerMap[key]; ok {
		return true
	} else {
		return false
	}
}

func (m KeyStringMap) Keys() *Slice {
	ret := make([]any, 0)
	for k, _ := range m.innerMap {
		ret = append(ret, k)
	}
	return ConvSlice(ret)
}

func (m KeyStringMap) Values() *Slice {
	ret := make([]any, 0)
	for _, v := range m.innerMap {
		ret = append(ret, v)
	}
	return ConvSlice(ret)
}

func (m KeyStringMap) EntrySet() []KSMapEntry {
	ret := make([]KSMapEntry, 0)
	for k, v := range m.innerMap {
		entry := KSMapEntry{Key: k, Value: v}
		ret = append(ret, entry)
	}
	return ret
}

func (m KeyStringMap) ContainsValue(value any) bool {
	values := m.Values()
	return slices.Contains(values.innerSlice, value)
}

func NewKSMap() *KeyStringMap {
	m := make(map[string]any)
	iMap := &KeyStringMap{innerMap: m}
	return iMap
}

func ConvKSMap(m any) *KeyStringMap {
	iMap := &KeyStringMap{innerMap: m.(map[string]any)}
	return iMap
}

func ConvJson2Map(jsonBytes []byte) (*KeyStringMap, error) {
	m := NewKSMap()
	err := json.Unmarshal(jsonBytes, &m.innerMap)
	return m, err
}
