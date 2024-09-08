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

func Int(number int) string {
	return strconv.Itoa(number)
}
