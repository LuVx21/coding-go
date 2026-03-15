package main

import (
	"fmt"
	"testing"

	"github.com/bits-and-blooms/bloom/v3"
)

func Test_bloom_00(t *testing.T) {
	filter := bloom.NewWithEstimates(1000000, 0.01)
	filter.Add([]byte("Love"))
	if filter.Test([]byte("Love")) {
		fmt.Println(true)
	}
}
