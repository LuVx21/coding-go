package container

import (
	"container/ring"
	"fmt"
	ringx "github.com/smallnest/exp/container/ring"
	"testing"
)

func Test_ring_00(t *testing.T) {
	r := ring.New(16)
	for i := 1; i <= 16; i++ {
		r.Value = i
		r = r.Next()
	}

	i, n := 0, r.Len()
	for p := r; i < n; p = p.Next() {
		fmt.Printf("%4d: %v <- %v <- %v\n", i, p.Prev().Value, p.Value, p.Next().Value)
		i++
	}
}

func Test_ring_01(t *testing.T) {
	r := ringx.New[int](16)
	for i := 1; i <= 16; i++ {
		r.Value = i
		r = r.Next()
	}
}
