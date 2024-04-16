package main

import (
    "fmt"
    "luvx.org/go/coding-common/common"
    "reflect"
)

func main() {
    p := common.NewPair("foo", "bar")
    v := reflect.ValueOf(&p)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    t := v.Type()
    for i := 0; i < v.NumField(); i++ {
        vField := v.Field(i)
        fmt.Println(vField.Interface())

        tField := t.Field(i)
        fmt.Printf("%s: %s\n", tField.Name, tField.Tag.Get("json"))
    }
}
