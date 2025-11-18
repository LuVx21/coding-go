package structs

import (
	"fmt"
	"reflect"

	"github.com/luvx21/coding-go/coding-common/common_x/pairs"
)

// ToMap 结构体转为Map[string]any
func ToMap(in any, tagName string) (map[string]any, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	} else if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("只能为结构体或其指针; 类型: %T", v)
	}

	result := make(map[string]any)
	t := v.Type()
	// 遍历结构体字段
	for i := range v.NumField() {
		vf := v.Field(i)
		field := t.Field(i)
		if vf.Kind() == reflect.Pointer {
			vf = vf.Elem()
		}
		if tagValue := field.Tag.Get(tagName); tagValue != "" {
			a := vf.Interface()
			switch vf.Kind() {
			case reflect.Struct, reflect.Pointer:
				a, _ = ToMap(a, tagName)
			default:
			}
			result[tagValue] = a
		}
	}
	return result, nil
}

// ToSingleMap 将结构体转为单层map
func ToSingleMap(in any, tag string) (map[string]any, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Pointer {
		// v = v.Elem()
	} else if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("只能为结构体或其指针; 类型: %T", v)
	}

	result := make(map[string]any)
	queue := make([]pairs.Pair[any, string], 0, 1)
	queue = append(queue, pairs.NewPair(in, ""))

	for len(queue) > 0 {
		e := queue[0]
		v := reflect.ValueOf(e.K)
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}
		queue = queue[1:]
		t := v.Type()
		for i := range v.NumField() {
			field := v.Field(i)
			kind := field.Kind()
			ti := t.Field(i)
			tagName := ti.Tag.Get(tag)
			if kind == reflect.Pointer {
				field = field.Elem()
				if field.Kind() == reflect.Struct {
					queue = append(queue, pairs.NewPair(field.Interface(), tagName+"."))
					continue
				}
			} else if kind == reflect.Struct {
				queue = append(queue, pairs.NewPair(field.Interface(), tagName+"."))
				continue
			}
			if tagName != "" {
				result[e.V+tagName] = field.Interface()
			}
		}
	}
	return result, nil
}
