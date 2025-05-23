package jsonpicker

import (
	"fmt"
	"reflect"
	"slices"
)

type Slice struct {
	innerSlice []any
}

func (rs Slice) ForEach(cb func(p any)) {
	for _, r := range rs.innerSlice {
		cb(r)
	}
}

func (rs Slice) Len() int {
	return len(rs.innerSlice)
}

func (rs Slice) MapToStrList(cb func(p any) string) []string {
	res := make([]string, 0)
	for _, r := range rs.innerSlice {
		sr := cb(r)
		res = append(res, sr)
	}
	return res
}

func (rs Slice) Map(cb func(p any) any) *Slice {
	res := NewSlice()
	for _, r := range rs.innerSlice {
		sr := cb(r)
		res.Append(sr)
	}
	return res
}

func (rs Slice) Filter(cb func(p any) bool) *Slice {
	res := NewSlice()
	for _, r := range rs.innerSlice {
		if ok := cb(r); ok {
			res.Append(r)
		}
	}
	return res
}

func (rs Slice) FindIndex(cb func(p any) bool) int {
	for idx, r := range rs.innerSlice {
		if ok := cb(r); ok {
			return idx
		}
	}
	return -1
}

func (rs Slice) IndexOf(p any) int {
	for idx, r := range rs.innerSlice {
		if r == p {
			return idx
		}
	}
	return -1
}

func (rs Slice) GetAt(index int) any {
	for idx, r := range rs.innerSlice {
		if idx == index {
			return r
		}
	}
	return nil
}

func (rs Slice) Contains(p any) bool {
	return slices.Contains(rs.innerSlice, p)
}

func (rs *Slice) Append(v any) {
	rs.innerSlice = append(rs.innerSlice, v)
}

func (rs *Slice) Extend(il any) {
	switch reflect.TypeOf(il).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(il)
		for i := range s.Len() {
			rs.Append(s.Index(i).Interface())
		}
	default:
		fmt.Println("[WARNING]Unsupported Type:", reflect.TypeOf(il).Kind())
	}
}

func (rs *Slice) Remove(v any) {
	idx := rs.IndexOf(v)
	if idx != -1 {
		rs.innerSlice = slices.Delete(rs.innerSlice, idx, idx+1)
	}
}

func (rs *Slice) RemoveAll(v any) {
	for {
		idx := rs.IndexOf(v)
		if idx != -1 {
			rs.innerSlice = slices.Delete(rs.innerSlice, idx, idx+1)
		} else {
			break
		}
	}
}

func ConvSlice(il any) (ret *Slice) {
	ret = NewSlice()
	switch reflect.TypeOf(il).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(il)
		for i := range s.Len() {
			ret.Append(s.Index(i).Interface())
		}
	default:
		ret = nil
	}
	return
}

func NewSlice() *Slice {
	i := make([]any, 0)
	return &Slice{innerSlice: i}
}
