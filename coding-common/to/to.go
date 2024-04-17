package to

import "time"

func ParseStr2Time(s string) (time.Time, error) {
    return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}
