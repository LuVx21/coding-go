package main

import (
    "bufio"
    "fmt"
    . "github.com/luvx21/coding-go/coding-common/common_x/funcs"
    "io"
    "iter"
    "slices"
    "strings"
    "testing"
)

//Backward 返回值是个函数, 该函数无出参, 入参仍为一个函数
func Backward[Slice ~[]E, E string](s Slice) iter.Seq2[int, E] {
    return func(yield func(int, E) bool) {
        for i := len(s) - 1; i >= 0; i-- {
            upper := strings.ToUpper(string(s[i]))
            yield(i, E(upper))
        }
        return
    }
}

func Pairs[V any](seq iter.Seq[V]) iter.Seq2[V, V] {
    return func(yield func(V, V) bool) {
        next, stop := iter.Pull(seq)
        defer stop()

        for {
            v1, hasNext := next()
            if !hasNext {
                return // 序列结束
            }

            v2, ok2 := next()
            if !ok2 {
                // 序列中有奇数个元素，最后一个元素没有配对
                return // 序列结束
            }

            if !yield(v1, v2) {
                return // 如果 yield 返回 false，停止迭代
            }
        }
    }
}

// Lines 返回一个迭代器，用于逐行读取 io.Reader 的内容, 迭代器不能重复使用
func Lines(r io.Reader) func(func(string) bool) {
    scanner := bufio.NewScanner(r)
    return func(yield func(string) bool) {
        for scanner.Scan() {
            if !yield(scanner.Text()) {
                return
            }
        }
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
    var fn = func(fun BiPredicate[int, byte]) {
        for i := 0; i < 26; i++ {
            if !fun(i, byte('a'+i)) {
                return
            }
        }
    }

    for k, v := range fn {
        fmt.Printf("%d: %c\n", k, v)
    }

    fn(func(k int, v byte) bool {
        fmt.Printf("%d: %c\n", k, v)
        return true
    })
}

func Test_rangefunc_02(t *testing.T) {
    sl := []string{"hello", "world", "golang"}

    for s := range slices.Values(sl) {
        fmt.Printf("%s\n", s)
    }

    for i, s := range slices.All(sl) {
        fmt.Printf("%d : %s\n", i, s)
    }

    for i, s := range slices.Backward(sl) {
        fmt.Printf("%d : %s\n", i, s)
    }
}
