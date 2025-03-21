package api

import (
	"fmt"
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
