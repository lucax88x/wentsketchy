package testutils

import (
	"testing"
	"time"

	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/stretchr/testify/require"
)

func DateToUtcTime(t *testing.T, date string) time.Time {
	return toUtcTime(t, clock.Date, date)
}

func DateTimeToUtcTime(t *testing.T, date string) time.Time {
	return toUtcTime(t, clock.DateTime, date)
}

func toUtcTime(t *testing.T, layout string, date string) time.Time {
	time, err := time.ParseInLocation(layout, date, time.UTC)

	if err != nil {
		require.NoError(t, err, "could not parse time as date %s", clock.Date)
	}

	return time
}
