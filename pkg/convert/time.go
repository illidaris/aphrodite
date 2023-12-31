package convert

import (
	"time"

	"github.com/spf13/cast"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	NumberTimeFormat  = "20060102150405"
	DateTimeFormat    = "20060102"
)

// ToCST set east +8
func ToCST(t time.Time) time.Time {
	zone := time.FixedZone("CST", 8*3600)
	return t.In(zone)
}

// layout: 2006-01-02 15:04:05 、 20060102150405 、2006年01月02日 15:04:05
func FormUnixToString(t int64, layout string) string {
	if t == 0 {
		return ""
	}
	return time.Unix(t, 0).Format(layout)
}

// TimeFormat format
func TimeFormat(t time.Time) string {
	return t.Format(DefaultTimeFormat)
}

// TimeNumber format number
func TimeNumber(t time.Time) int64 {
	return cast.ToInt64(t.Format(NumberTimeFormat))
}

// TimeNumber format number date
func TimeDate(t time.Time) int64 {
	return cast.ToInt64(t.Format(DateTimeFormat))
}
