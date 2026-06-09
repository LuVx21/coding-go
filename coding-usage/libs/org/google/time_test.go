package x

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func Test_x_time_00(t *testing.T) {
	s := rate.Sometimes{
		First:    2,
		Every:    2,
		Interval: 2 * time.Second,
	}
	s.Do(func() { fmt.Println("1 (First:2)") })
	s.Do(func() { fmt.Println("2 (First:2)") })
	s.Do(func() { fmt.Println("3 (Every:2)") })
	time.Sleep(2 * time.Second)
	s.Do(func() { fmt.Println("4 (Interval)") })
	s.Do(func() { fmt.Println("5 (Every:2)") })
	s.Do(func() { fmt.Println("6") })
}
