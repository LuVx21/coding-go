package bili

import (
    "github.com/luvx21/coding-go/coding-common/common_x"
    "luvx/gin/service"
)

func RunnerRegister() []*service.Runner {
    return []*service.Runner{
        {Name: "拉取bili_up", Crontab: "17 17 9/12 * * *", Fn: func() { common_x.RunCatching(PullAllUpVideo) }},
    }
}
