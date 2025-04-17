package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestPass(t *testing.T) {
	rl := New(1, 3)

	for range 60 {
		res := rl.Pass()
		time.Sleep(100 * time.Millisecond)
		fmt.Println(res)
	}
}
