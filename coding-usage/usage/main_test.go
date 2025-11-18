package main

import "testing"

var x, y int
var ( // 这种因式分解关键字的写法一般用于声明全局变量
	a int
	b bool
)

var c, d int = 1, 2
var e, f = 123, "foo"

func Test_main_00(t *testing.T) {
	// 这种不带声明格式的只能在函数体中出现
	g, h := 456, "bar"
	println(x, y)
	println(a, b)
	println(c, d)
	println(e, f)
	println(g, h)
}
