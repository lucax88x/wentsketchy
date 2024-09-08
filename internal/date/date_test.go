package date_test

import (
	"testing"
	"time"

	"github.com/lucax88x/wentsketchy/internal/date"
	"github.com/lucax88x/wentsketchy/internal/formatter"
	"github.com/stretchr/testify/require"
)

func TestUnitShouldGetStartOfMonth(t *testing.T) {
	// GIVEN
	parsedDate, _ := time.Parse(date.Date, "2023-06-08")

	// WHEN
	result := date.StartOfMonth(parsedDate)

	// THEN
	require.Equal(t, "2023-06-01 00:00:00", formatter.DateTime(result))
}

func TestUnitShouldGetEndOfMonth(t *testing.T) {
	// GIVEN
	parsedDate, _ := time.Parse(date.Date, "2023-06-08")

	// WHEN
	result := date.EndOfMonth(parsedDate)

	// THEN
	require.Equal(t, "2023-06-30 23:59:59", formatter.DateTime(result))
}
