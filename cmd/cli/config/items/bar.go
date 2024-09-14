package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Bar(batches [][]string) ([][]string, error) {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   45,
		// TODO: move to settings
		ColorOptions: sketchybar.ColorOptions{
			Color: settings.ColorBarColor,
		},
		BorderOptions: sketchybar.BorderOptions{
			Width: 2,
			Color: settings.ColorBarBorderColor,
		},
		Shadow: false,
		Sticky: true,
		PaddingOptions: sketchybar.PaddingOptions{
			Right: 10,
			Left:  10,
		},
		YOffset: -5,
		Margin:  -2,
		Topmost: "window",
	}

	batches = batch(batches, m(s("--bar"), bar.ToArgs()))

	return batches, nil
}
