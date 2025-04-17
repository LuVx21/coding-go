package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func Test_chromedp_00(t *testing.T) {
	edgePath := "/Applications/Microsoft Edge Dev.app/Contents/MacOS/Microsoft Edge Dev"

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(edgePath),
		// chromedp.Flag("remote-debugging-port", "9222"), // 指定调试端口
		// chromedp.UserDataDir("C:\\path\\to\\edge\\userdata"),                                   // 可选：指定用户数据目录
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 创建 chromedp 上下文
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.baidu.com"),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML("html", &html),
	)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(html)
}
