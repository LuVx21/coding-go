package fmt_x

import (
	"fmt"
	"testing"
	"time"
)

func Test_00(t *testing.T) {
}

func Test_01(t *testing.T) {
	PrintlnRow(1, "aa", time.Now())
	Println([]any{"aa", 2}, []any{"a", 1}, []any{2, time.Now()})
}

func Test_02(t *testing.T) {
	r1 := []any{"a", 22, []string{"a", "bb", "ccc"}, map[string]any{"a": true, "b": "bbb"}}
	r2 := []any{"b", 33, []string{"b", "bb", "ccc"}, map[string]any{"a": true, "b": "bbb"}, struct{ string }{"ss"}}
	Println0(r1, r2)

	PrintlnRow0("a", 22, []string{"a", "bb", "ccc"}, "ʕ◔ϖ◔ʔ")
}

func Test_color_00(t *testing.T) {
	// 31: 红色前景
	fmt.Println(colorRed + "这是红色文本" + colorReset)

	// 32: 绿色前景
	fmt.Println(colorGreen + "这是绿色文本" + colorReset)

	// 33: 黄色前景
	fmt.Println(colorYellow + "这是黄色文本" + colorReset)

	// 带背景色
	fmt.Println(a + "红色文本绿色背景" + colorReset)

	// 加粗
	fmt.Println(b + "加粗红色文本" + colorReset)
}
