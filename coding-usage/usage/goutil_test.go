package main

import (
	"fmt"
	"testing"

	"github.com/Goldziher/go-utils/maputils"
	"github.com/gookit/goutil"
	. "github.com/luvx21/coding-go/coding-common/maps_x"
)

func Test_01(t *testing.T) {
	empty := goutil.IsEmpty(nil)
	fmt.Println(empty)
}

func Test_11_00(t *testing.T) {
	m := Map[string, bool]{"a": true}
	keys := maputils.Keys(m)
	fmt.Println(keys)
}
