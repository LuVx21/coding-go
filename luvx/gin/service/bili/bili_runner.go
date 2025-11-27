package bili

import (
	"luvx/gin/service"

	"github.com/luvx21/coding-go/coding-common/common_x"
)

func RunnerRegister() []*service.Runner {
	return []*service.Runner{
		{Name: "拉取bili_up", Crontab: "17 17 9/12 * * *", Fn: func() { common_x.RunCatching(PullAllUpVideo) }},
		// {Name: "拉取bili_flow", Crontab: "17 17 3/8 * * *", Fn: func() {
		// 	common_x.RunCatching(func() { service.RunnerLocker.LockRun("拉取bili_flow", time.Minute*10, func() { timeFlow() }) })
		// }},
	}
}
