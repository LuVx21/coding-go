package exp

import (
	"testing"

	"github.com/puzpuzpuz/xsync/v3"
)

func Test_map_00(t *testing.T) {
	m := xsync.NewMap()
	m.Store("foo", "bar")
	println(m.Load("foo"))
	println(m.Size())
}
