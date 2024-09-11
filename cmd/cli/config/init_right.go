package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func initRight(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	batches, err := initCalendar()

	if err != nil {
		return fmt.Errorf("init calendar %w", err)
	}
	err = di.Sketchybar.Run(ctx, batches)

	if err != nil {
		return fmt.Errorf("apply to sketchybar %w", err)
	}

	return nil
}

func initCalendar() ([]string, error) {
	var batches = make([][]string, 0)

	calendar := sketchybar.ItemOptions{
		Icon: sketchybar.ItemIconOptions{
			Value: wentsketchy.IconClock,
			Font: sketchybar.FontOptions{
				Font: wentsketchy.FontIcon,
				Kind: "Regular",
				Size: "12.0",
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
		},
		Background: sketchybar.BackgroundOptions{
			BorderOptions: sketchybar.BorderOptions{
				Width: 2,
				Color: wentsketchy.ColorBackground1,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: wentsketchy.ColorBackground1,
			},
		},
		UpdateFreq: 30,
		Script:     "wentsketchy update calendar",
		// Click_script:            "$PLUGIN_DIR/zen.sh",
	}

	batches = batch(batches, s("--add", "item", "calendar", "right"))
	batches = batch(batches, m(s("--set", "calendar"), calendar.ToArgs()))
	batches = batch(batches, s("--subscribe", "calendar", "system_woke"))

	return flatten(batches...), nil
}
