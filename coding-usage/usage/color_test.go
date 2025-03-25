package main

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	color2 "github.com/gookit/color"
)

func Test_color_00(t *testing.T) {
	// 格式: \033[显示方式;前景色;背景色m
	// 31: 红色前景
	fmt.Println("\033[31m这是红色文本\033[0m")

	// 32: 绿色前景
	fmt.Println("\033[32m这是绿色文本\033[0m")

	// 33: 黄色前景
	fmt.Println("\033[33m这是黄色文本\033[0m")

	// 带背景色
	fmt.Println("\033[31;42m红色文本绿色背景\033[0m")

	// 加粗
	fmt.Println("\033[1;31m加粗红色文本\033[0m")
}

func Test_color_01(t *testing.T) {
	color.Red("这是红色文本")
	color.Green("这是绿色文本")
	color.Yellow("这是黄色文本")

	// 自定义样式
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("青色带下划线的文本")

	// 组合颜色
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	boldRed.Println("加粗的红色文本")
}
func Test_color_02(t *testing.T) {
	color2.Red.Println("红色文本")
	color2.Green.Print("绿色文本")
	color2.Yellow.Printf("黄色文本 %s\n", "带格式化")

	// 自定义样式
	style := color2.New(color2.FgWhite, color2.BgBlack, color2.OpBold)
	style.Println("自定义样式文本")
}
