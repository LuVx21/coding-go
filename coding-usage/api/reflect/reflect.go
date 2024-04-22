package main

import (
    "fmt"
    "github.com/luvx21/coding-go/coding-common/common_x"
    "reflect"
)

func main() {
    p := common_x.NewPair("foo", "bar")
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
