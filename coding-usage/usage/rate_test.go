package main

import (
	"context"
	"fmt"
	"github.com/luvx21/coding-go/coding-common/times_x"
	"golang.org/x/time/rate"
	"testing"
)

func Test_rate_00(t *testing.T) {
	limiter := rate.NewLimiter(1, 1)
	for i := range 10 {
		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Println("Error waiting for limiter:", err)
			return
		}
		fmt.Println(times_x.TimeNow(), i+1)
	}
}
