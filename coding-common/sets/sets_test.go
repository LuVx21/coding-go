package sets

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	s := Set[string]{"foobar": struct{}{}}
	s.Add("foo", "bar")
	fmt.Println(s.Contains("foo"), s.Contains("bar1"))
}
