package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-viper/mapstructure/v2"
)

func Test_mapstructure_00(t *testing.T) {
	type user struct {
		Name   string
		Age    int
		Emails []string
	}

	var u user
	err := mapstructure.Decode(map[string]any{
		"name":   "foo",
		"Age":    18,
		"emails": []string{"aaa", "bbb"},
	}, &u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)
}
