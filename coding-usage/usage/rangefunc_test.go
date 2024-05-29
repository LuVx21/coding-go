package main

import (
    "fmt"
    "testing"
)

func Backward[E any](s []E) func(func(int, E) bool) {
    return func(yield func(int, E) bool) {
        for i := len(s) - 1; i >= 0; i-- {
            if !yield(i, s[i]) {
                return
            }
        }
        return
    }
}

//export GOEXPERIMENT=rangefunc
func Test_rangefunc_00(t *testing.T) {
    sl := []string{"hello", "world", "golang"}
    for i, s := range Backward(sl) {
        fmt.Printf("%d : %s\n", i, s)
    }
}

func Test_rangefunc_01(t *testing.T) {
    var fn = func(fun func(k int, v byte) bool) {
        for i := 0; i < 26; i++ {
            if !fun(i, byte('a'+i)) {
                return
            }
        }
    }

    for k, v := range fn {
        fmt.Printf("%d: %c\n", k, v)
    }
}

func Test_rangefunc_02(t *testing.T) {
    //sl := []string{"hello", "world", "golang"}
    //
    //for i, s := range slices.All(sl) {
    //    fmt.Printf("%d : %s\n", i, s)
    //}
    //
    //for i, s := range slices.Backward(sl) {
    //    fmt.Printf("%d : %s\n", i, s)
    //}
}
