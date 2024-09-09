package formatter

import (
	"strconv"
	"time"

	"github.com/lucax88x/wentsketchy/internal/clock"
)

func DateTime(time time.Time) string {
	return time.Format(clock.DateTime)
}

func Date(time time.Time) string {
	return time.Format(clock.Date)
}

func Time(time time.Time) string {
	return time.Format(clock.Time)
}

func HoursMinutes(time time.Time) string {
	return time.Format(clock.HoursMinutes)
}

func Int(number int) string {
	return strconv.Itoa(number)
}
