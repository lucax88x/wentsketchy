package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Bar(batches [][]string) ([][]string, error) {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   settings.SketchybarSettings.BarHeight,
		Margin:   settings.SketchybarSettings.BarMargin,
		YOffset:  0,
		Padding: sketchybar.PaddingOptions{
			Right: 8,
			Left:  8,
		},
		Topmost:       "off",
		Sticky:        "on",
		Shadow:        "off",
		FontSmoothing: "on",
		Color: sketchybar.ColorOptions{
			Color: settings.SketchybarSettings.BarBackgroundColor,
		},
	}

	batches = batch(batches, m(s("--bar"), bar.ToArgs()))

	return batches, nil
}
