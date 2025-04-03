package main

import (
	"testing"

	"github.com/fatih/color"
	color2 "github.com/gookit/color"
)

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
