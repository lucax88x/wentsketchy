package config

import (
	"context"
	"fmt"
	"time"

	"github.com/lucax88x/wentsketchy/internal/formatter"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func Update(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	err := updateCalendar(ctx, di)

	if err != nil {
		return fmt.Errorf("update calendar %w", err)
	}

	return nil
}

func updateCalendar(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	calendar := sketchybar.ItemOptions{
		Label: sketchybar.ItemLabelOptions{
			Value: formatter.HoursMinutes(time.Now()),
		},
	}

	return di.Sketchybar.Run(
		ctx,
		m(s("--set", "calendar"), calendar.ToArgs()),
	)
}
