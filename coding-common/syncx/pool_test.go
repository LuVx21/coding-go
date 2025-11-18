package syncx

import (
	"bytes"
	"fmt"
	"testing"
)

var p = NewPool(func() *bytes.Buffer {
	return new(bytes.Buffer)
})

func Test_pool_00(t *testing.T) {
	buf := p.Get()
	defer p.Put(buf)
	buf.WriteString("hello")

	fmt.Println(buf.String())
}

// 不使用 Pool
func Benchmark_WithoutPool(b *testing.B) {
	for b.Loop() {
		buf := new(bytes.Buffer)
		buf.WriteString("test")
		_ = buf.String()
	}
}

// 使用 Pool
func Benchmark_WithPool(b *testing.B) {
	for b.Loop() {
		buf := p.Get()
		buf.WriteString("test")
		s := buf.String()
		buf.Reset()
		p.Put(buf)
		_ = s
	}
}
