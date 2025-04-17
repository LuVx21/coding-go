package api

import (
	"fmt"
	"testing"
	"time"
)

func Test_time_00(t *testing.T) {
	time.AfterFunc(time.Second*5, func() { fmt.Println("定时") })
	time.Sleep(time.Second * 6)
}
