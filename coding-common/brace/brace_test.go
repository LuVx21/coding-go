package brace

import (
	"fmt"
	"testing"
)

func Test_Expand_00(t *testing.T) {
	ss := []string{
		"a{b,c}d",
		"x{a,{b{1..3},c{,d,e}}}y",
		"file_{1..3}.txt",
		"{A..c..2}",
		"{11..19..3}",
		"{-15..1}",
		`a\{b,c\}`,
		"a{}b",
	}
	for _, s := range ss {
		fmt.Println(s, " -> ", Expand(s))
	}
}

func Test_Expand_01(t *testing.T) {
}
