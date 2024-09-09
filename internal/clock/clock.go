package clock

import "time"

type Clock interface {
	Now() time.Time
}

const Date = "2006-01-02"
const DateTime = "2006-01-02 15:04:05"
const Time = "15:04:05"
const HoursMinutes = "15:04"

type SystemCock struct{}

func NewSystemCock() Clock {
	return &SystemCock{}
}

func (r *SystemCock) Now() time.Time {
	return time.Now()
}
