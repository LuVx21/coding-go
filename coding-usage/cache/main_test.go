package main

import (
    "fmt"
    "github.com/patrickmn/go-cache"
    "testing"
    "time"
)

func Test01(t *testing.T) {
    c := cache.New(5*time.Minute, 10*time.Minute)
    c.Set("foo", "bar", cache.DefaultExpiration)
    v, exist := c.Get("foo")
    fmt.Println(v, exist)
}
