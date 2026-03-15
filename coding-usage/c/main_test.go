package main

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/ebitengine/purego"
)

func Test_main_00(in *testing.T) {
	var libPath string

	// 根据操作系统确定库路径
	switch runtime.GOOS {
	case "darwin":
		libPath = "./lib/libadd.dylib"
	case "linux":
		libPath = "./lib/libadd.so"
	case "windows":
		libPath = "./lib/add.dll"
	default:
		panic("unsupported platform")
	}

	fmt.Println("lib路径", libPath)
	lib, err := purego.Dlopen(libPath, purego.RTLD_NOW)
	if err != nil {
		panic(err)
	}

	// 2. 定义函数变量（注意参数类型需与C一致）
	var add func(int, int) int
	// 3. 注册库函数
	purego.RegisterLibFunc(&add, lib, "add")

	// 4. 调用函数
	result := add(30, 40)
	fmt.Printf("30 + 40 = %d\n", result)
}
