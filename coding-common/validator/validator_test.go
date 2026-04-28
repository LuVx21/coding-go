package validator

import (
	"fmt"
	"testing"
)

func Test_validator_00(t *testing.T) {
	m := make(map[string]any, 0)

	IfThen(true, func() { m["a"] = "aa" })

	IfThenI("a", func(s string) bool { return len(s) > 0 }, func(s string) { m["b"] = "bb" })

	fmt.Println(m)
}
