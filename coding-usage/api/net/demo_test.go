package main

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_http(t *testing.T) {
	r, _ := http.Get("https://httpbin.org/get")
	fmt.Println(r.StatusCode)
}

func Test_a(t *testing.T) {
	main()
}
