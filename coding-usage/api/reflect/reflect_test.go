package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/luvx21/coding-go/coding-common/fmt_x"
)

func aa(s any) []any {
	result := make([]any, 0)

	ts := reflect.TypeOf(s)
	result = append(result, ts, ts.Kind())
	if ts.Kind() == reflect.Pointer {
		result = append(result, ts.Elem(), ts.Elem().Kind())
	} else {
		result = append(result, "无", "无")
	}

	vs := reflect.ValueOf(s)
	result = append(result, vs, vs.Kind(), vs.Interface())
	if vs.Kind() == reflect.Pointer {
		result = append(result, vs.Elem(), vs.Elem().Kind(), vs.Elem().Interface())
	} else {
		result = append(result, "无", "无", "无")
	}
	return result
}

func Test_00(t *testing.T) {
	var s1 json.Number = "hello"
	var s string = "hello"
	var p1 = &s
	var p2 = p1
	var p3 = p2

	fmt_x.Println(
		[]any{"Type", "Type Kind", "Type Elem", "Type Elem Kind", "Value", "Value Kind", "Value Interface", "Value Elem", "Value Elem Kind", "Value Elem Interface"},
		aa(s1),
		aa(&s1),
		aa(s),
		aa(p1),
		aa(p2),
		aa(p3),
	)
}
