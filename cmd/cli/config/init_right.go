package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func initRight(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	err := initCalendar(ctx, di)

	if err != nil {
		return fmt.Errorf("init calendar %w", err)
	}

	return nil
}

func initCalendar(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	calendar := sketchybar.ItemOptions{
		Icon: sketchybar.ItemIconOptions{
			Value: IconClock,
			Font: sketchybar.FontOptions{
				Font: iconFont,
				Kind: "Regular",
				Size: "12.0",
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  10,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  10,
			},
		},
		Background: sketchybar.BackgroundOptions{
			BorderOptions: sketchybar.BorderOptions{
				Width: 2,
				Color: ColorBackground1,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: ColorBackground1,
			},
		},
		UpdateFreq: 30,
		Script:     "wentsketchy update calendar",
		// Click_script:            "$PLUGIN_DIR/zen.sh",
	}

	return di.Sketchybar.Run(
		ctx,
		s("--add", "item", "calendar", "right"),
		m(s("--set", "calendar"), calendar.ToArgs()),
		s("--subscribe", "calendar", "system_woke"),
	)
}
