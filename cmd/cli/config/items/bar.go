package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Bar(batches Batches) (Batches, error) {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   settings.Sketchybar.BarHeight,
		Margin:   settings.Sketchybar.BarMargin,
		YOffset:  pointer(-40),
		Padding: sketchybar.PaddingOptions{
			Right: pointer(8),
			Left:  pointer(8),
		},
		Topmost:       "off",
		Sticky:        "on",
		Shadow:        "off",
		FontSmoothing: "on",
		Color: sketchybar.ColorOptions{
			Color: settings.Sketchybar.BarBackgroundColor,
		},
	}

	batches = batch(batches, m(s("--bar"), bar.ToArgs()))

	return batches, nil
}

func ShowBar(batches Batches) (Batches, error) {
	bar := sketchybar.BarOptions{
		YOffset: pointer(0),
	}

	batches = batch(batches, m(s(
		"--animate",
		sketchybar.AnimationTanh,
		settings.Sketchybar.BarTransitionTime,
		"--bar",
	), bar.ToArgs()))

	return batches, nil
}
