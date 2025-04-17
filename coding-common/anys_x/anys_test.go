package anys_x

import (
	"testing"

	"github.com/luvx21/coding-go/coding-common/test"
)

func Test_any_00(t *testing.T) {
	tests := []test.Step{
		{Name: "用例1", Input: 11, Expected: "11"},
		{Name: "用例2", Input: "abc", Expected: "abc"},
	}

	test.OneOne(t, tests, func(v any) any {
		return String(v)
	})
}

func Test_any_01(t *testing.T) {
}
