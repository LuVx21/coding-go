package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/charmbracelet/glamour"
	"github.com/yuin/goldmark"
)

var md = `# Hello World

This is a simple example of Markdown rendering with Glamour!
Check out the [other examples](https://github.com/charmbracelet/glamour/tree/master/examples) too.

    select *
    from user;

# title

> 注释

* 1
* 2
  - 22

---

### 111(deepseek-chat) 2025-03-29 13:26:50
你好！我是 **DeepSeek Chat**，由深度求索（DeepSeek）公司开发的智能AI助手。我可以帮助你解答各种问题，包括学习、工作、编程、写作、翻译、生活建议等。

### **我的特点：**
✅ **免费使用**：目前无需付费，随时可以向我提问！
✅ **知识丰富**：我的知识截止到 **2024年7月**，能提供较新的信息。
✅ **超长上下文**：支持 **128K** 上下文，可以处理超长文档和复杂对话。
✅ **文件阅读**：可以上传 **PDF、Word、Excel、PPT、TXT** 等文件，帮你提取和分析内容。
✅ **编程助手**：支持多种编程语言，能帮你调试代码、优化算法、解释技术概念。
✅ **写作与创意**：可以帮你写文章、改写文案、生成故事、头脑风暴创意点子。

### **我不能做什么？**
❌ 无法实时联网（但你可以手动上传最新资料让我分析）。
❌ 不能生成图片、音频或视频（但可以帮你写脚本或描述）。

无论是学习、工作还是日常问题，都可以来找我聊聊！😊 你今天想了解什么呢？

| a    | b    | c    |
| :--- | :--- | :--- |
| aa   | bb   | cc   |

Bye!
`

func Test_glamour(t *testing.T) {
	var render, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
		// glamour.WithStyles(styles.DraculaStyleConfig),
	)

	out, _ := render.Render(md)
	// out, _ = glamour.Render(md, "light")
	fmt.Println(out)
}

func Test_goldmark(t *testing.T) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
}
