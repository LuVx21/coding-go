package retry

import (
	"fmt"
	"log"
	"time"

	"github.com/luvx21/coding-go/coding-common/strings_x"
	"github.com/luvx21/coding-go/coding-common/times_x"
)

// SupplyWithRetry fn: 内部异常或错误, 使用panic
func SupplyWithRetry[T any](
	name string, fn func() T,
	maxRetryTimes int32, retryPeriod time.Duration,
) (T, error) {
	if strings_x.IsBlank(name) {
		name = "重试" + times_x.TimeNow()
	}

	var times int32
	for times <= maxRetryTimes {
		var pass = true
		t := func() T {
			defer func() {
				if r := recover(); r != nil {
					pass = false
					time.Sleep(retryPeriod)
					times++
					if times <= maxRetryTimes {
						log.Printf("%s->当前异常:%v, 进行第%d次重试", name, r, times)
					}
				}
			}()
			return fn()
		}()
		if pass {
			return t, nil
		}
	}
	log.Printf("进行%d次重试仍失败", maxRetryTimes)
	var zero T
	return zero, fmt.Errorf("重试失败")
}
