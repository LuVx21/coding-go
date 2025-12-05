package sets

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	s := set[string]{"foobar": struct{}{}}
	s.Add("foo", "bar")
	fmt.Println(s.Contain("foo"), s.Contain("bar1"))
}
