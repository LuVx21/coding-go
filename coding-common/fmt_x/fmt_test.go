package fmt_x

import (
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
