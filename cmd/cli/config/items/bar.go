package items

import (
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Bar(batches [][]string) ([][]string, error) {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   45,
		Shadow:   false,
		Sticky:   true,
		PaddingOptions: sketchybar.PaddingOptions{
			Right: 10,
			Left:  10,
		},
		YOffset: -5,
		Topmost: "window",
	}

	batches = batch(batches, m(s("--bar"), bar.ToArgs()))

	return batches, nil
}
