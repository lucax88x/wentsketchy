package date

import "time"

const Date = "2006-01-02"
const DateTime = "2006-01-02 15:04:05"

func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}
