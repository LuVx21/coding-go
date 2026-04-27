package time

import (
	"fmt"
	"testing"
	"time"

	"github.com/luvx21/coding-go/coding-common/times_x"
)

func Test_time_ticker_00(t *testing.T) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println(times_x.TimeNowSecond())
	}
}
