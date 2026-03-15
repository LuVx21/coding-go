package main

import (
	"fmt"
	"testing"

	setx "github.com/smallnest/exp/container/set"
)

func Test_set_01(t *testing.T) {
	s := setx.NewSet[string]()
	s.Add("a", "b")
	fmt.Println(s)
}
