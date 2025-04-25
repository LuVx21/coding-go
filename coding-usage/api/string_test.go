package api

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_xx(t *testing.T) {
	str := "Hello, 世界"

	for i := range len(str) {
		fmt.Printf("Byte %d: %x\n", i, str[i])
	}
	for i, char := range str {
		fmt.Printf("Character %d: %c\n", i, char)
	}

	fmt.Println([]byte(str), []rune(str))

	var chars [13]byte
	copy(chars[:], str)
	fmt.Println("Character array:", chars)
}

func BenchmarkStringBuilder(b *testing.B) {
	var sb strings.Builder
	for b.Loop() {
		sb.WriteString("hello")
	}
	_ = sb.String()
}

func BenchmarkBytesBuffer(b *testing.B) {
	var buf bytes.Buffer
	for b.Loop() {
		buf.WriteString("hello")
	}
	_ = buf.String()
}
