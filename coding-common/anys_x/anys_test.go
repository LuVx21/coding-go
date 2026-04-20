package anys_x

import (
	"testing"

	"github.com/luvx21/coding-go/coding-common/tests"
)

func Test_any_00(t *testing.T) {
	cases := []tests.Step{
		{Name: "用例1", Input: 11, Expected: "11"},
		{Name: "用例2", Input: "abc", Expected: "abc"},
	}

	tests.OneOne(t, cases, func(v any) any {
		return String(v)
	})
}

func Test_any_01(t *testing.T) {
}
