package times_x

import "time"

func ParseDateTimeLocal(timeStr string) (time.Time, error) {
	return ParseDateTime(timeStr, time.Local)
}
func ParseDateTime(timeStr string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(time.DateTime, timeStr, loc)
}
