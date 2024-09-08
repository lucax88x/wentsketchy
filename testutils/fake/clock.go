package fake

import (
	"time"

	"github.com/lucax88x/wentsketchy/internal/clock"
)

type Clock struct {
	Time time.Time
}

func (m *Clock) Now() time.Time {
	return m.Time
}

var _ clock.Clock = (*Clock)(nil)
