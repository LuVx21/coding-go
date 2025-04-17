package jsonpicker

import "slices"

type Map struct {
	innerMap map[any]any
}

type MapEntry struct {
	Key   any
	Value any
}

func (m *Map) Get(key any) any {
	if v, ok := m.innerMap[key]; ok {
		return v
	} else {
		return nil
	}
}

func (m *Map) Set(key, value any) {
	m.innerMap[key] = value
}

func (m *Map) ContainsKey(key any) bool {
	if _, ok := m.innerMap[key]; ok {
		return true
	} else {
		return false
	}
}

func (m *Map) Keys() *Slice {
	ret := make([]any, 0)
	for k, _ := range m.innerMap {
		ret = append(ret, k)
	}
	return ConvSlice(ret)
}

func (m *Map) Values() *Slice {
	ret := make([]any, 0)
	for _, v := range m.innerMap {
		ret = append(ret, v)
	}
	return ConvSlice(ret)
}

func (m *Map) EntrySet() []MapEntry {
	ret := make([]MapEntry, 0)
	for k, v := range m.innerMap {
		entry := MapEntry{Key: k, Value: v}
		ret = append(ret, entry)
	}
	return ret
}

func (m *Map) ContainsValue(value any) bool {
	values := m.Values()
	return slices.Contains(values.innerSlice, value)
}

func NewMap() *Map {
	m := make(map[any]any)
	iMap := &Map{innerMap: m}
	return iMap
}
