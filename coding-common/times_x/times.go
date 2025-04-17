package times_x

import (
	"fmt"
	"time"
)

type TimeFormatType int
type TimeFormat struct {
	Format string
	Typ    TimeFormatType
}

const (
	Day       = time.Hour * 24
	Year      = 365 * Day
	YEAR_LEAP = 366 * Day

	TimeFormatNoTimezone TimeFormatType = iota
	TimeFormatNamedTimezone
	TimeFormatNumericTimezone
	TimeFormatNumericAndNamedTimezone
	TimeFormatTimeOnly
)

var (
	TimeFormats = []TimeFormat{
		{"02 Jan 2006", TimeFormatNoTimezone},
		{"2006-01-02 15:04:05 -07:00", TimeFormatNumericTimezone},
		{"2006-01-02 15:04:05 -0700", TimeFormatNumericTimezone},
		{"2006-01-02 15:04:05.999999999 -0700 MST", TimeFormatNumericAndNamedTimezone},
		{"2006-01-02 15:04:05Z07:00", TimeFormatNumericTimezone},
		{"2006-01-02 15:04:05Z0700", TimeFormatNumericTimezone},
		{"2006-01-02T15:04:05", TimeFormatNoTimezone},
		{"2006-01-02T15:04:05-0700", TimeFormatNumericTimezone},
		{time.Layout, TimeFormatNumericTimezone},
		{time.ANSIC, TimeFormatNoTimezone},
		{time.UnixDate, TimeFormatNamedTimezone},
		{time.RubyDate, TimeFormatNumericTimezone},
		{time.RFC822, TimeFormatNamedTimezone},
		{time.RFC822Z, TimeFormatNumericTimezone},
		{time.RFC850, TimeFormatNamedTimezone},
		{time.RFC1123, TimeFormatNamedTimezone},
		{time.RFC1123Z, TimeFormatNumericTimezone},
		{time.RFC3339, TimeFormatNumericTimezone},
		{time.RFC3339Nano, TimeFormatNumericTimezone},
		{time.Kitchen, TimeFormatTimeOnly},

		{time.Stamp, TimeFormatTimeOnly},
		{time.StampMilli, TimeFormatTimeOnly},
		{time.StampMicro, TimeFormatTimeOnly},
		{time.StampNano, TimeFormatTimeOnly},
		{time.DateTime, TimeFormatNoTimezone},
		{time.DateOnly, TimeFormatNoTimezone},
		{time.TimeOnly, TimeFormatTimeOnly},
	}
)

func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, time.UTC, TimeFormats)
}

func StringToDateInDefaultLocation(s string, location *time.Location) (time.Time, error) {
	return parseDateWith(s, location, TimeFormats)
}

func parseDateWith(s string, location *time.Location, formats []TimeFormat) (d time.Time, e error) {
	for _, format := range formats {
		if d, e = time.Parse(format.Format, s); e == nil {
			if format.Typ == TimeFormatNamedTimezone || format.Typ == TimeFormatNoTimezone {
				if location == nil {
					location = time.Local
				}
				year, month, day := d.Date()
				hour, _min, sec := d.Clock()
				d = time.Date(year, month, day, hour, _min, sec, d.Nanosecond(), location)
			}
			return
		}
	}
	return d, fmt.Errorf("unable to parse date: %s", s)
}

func TimeNowDate() string {
	return time.Now().Format(time.DateOnly)
}

func TimeNowDateSecond() string {
	return time.Now().Format(time.DateTime)
}

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05.999999")
}

func ParseStr2Time(s string) (time.Time, error) {
	return time.ParseInLocation(time.DateTime, s, time.Local)
}
