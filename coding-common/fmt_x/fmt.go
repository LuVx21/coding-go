package fmt_x

import (
    "time"
)

func TimeNowDate() string {
    return time.Now().Format("2006-01-02")
}

func TimeNowDateSecond() string {
    return time.Now().Format("2006-01-02")
}

func TimeNow() string {
    return time.Now().Format("2006-01-02 15:04:05.999999")
}
