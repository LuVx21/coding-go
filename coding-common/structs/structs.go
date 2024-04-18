package structs

import (
    "fmt"
    "github.com/luvx21/coding-go/coding-common/common"
    "reflect"
)

// ToMap 结构体转为Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
    } else if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("只能为结构体或其指针; 类型: %T", v)
    }

    result := make(map[string]interface{})
    t := v.Type()
    // 遍历结构体字段
    for i := 0; i < v.NumField(); i++ {
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
func ToSingleMap(in interface{}, tag string) (map[string]interface{}, error) {
    v := reflect.ValueOf(in)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
    } else if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("只能为结构体或其指针; 类型: %T", v)
    }

    result := make(map[string]interface{})
    queue := make([]common.Pair[interface{}, string], 0, 1)
    queue = append(queue, common.NewPair(in, ""))

    for len(queue) > 0 {
        e := queue[0]
        v := reflect.ValueOf(e.K)
        if v.Kind() == reflect.Pointer {
            v = v.Elem()
        }
        queue = queue[1:]
        t := v.Type()
        for i := 0; i < v.NumField(); i++ {
            field := v.Field(i)
            kind := field.Kind()
            ti := t.Field(i)
            tagName := ti.Tag.Get(tag)
            if kind == reflect.Pointer {
                field = field.Elem()
                if field.Kind() == reflect.Struct {
                    queue = append(queue, common.NewPair(field.Interface(), tagName+"."))
                    continue
                }
            } else if kind == reflect.Struct {
                queue = append(queue, common.NewPair(field.Interface(), tagName+"."))
                continue
            }
            if tagName != "" {
                result[e.V+tagName] = field.Interface()
            }
        }
    }
    return result, nil
}
