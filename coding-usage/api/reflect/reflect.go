package main

import (
	"fmt"
	"reflect"

	"github.com/luvx21/coding-go/coding-common/common_x/pairs"
)

func main() {
	p := pairs.NewPair("foo", "bar")
	v := reflect.ValueOf(&p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := range v.NumField() {
		vField := v.Field(i)
		fmt.Println(vField.Interface())

		tField := t.Field(i)
		fmt.Printf("%s: %s\n", tField.Name, tField.Tag.Get("json"))
	}
}
