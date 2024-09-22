package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Defaults(batches [][]string) ([][]string, error) {
	defaults := sketchybar.ItemOptions{
		YOffset: 0,
		Padding: sketchybar.PaddingOptions{
			Right: 0,
			Left:  0,
		},
		Icon: sketchybar.ItemIconOptions{
			Color: sketchybar.ColorOptions{
				Color: settings.SketchybarSettings.IconColor,
			},
			Font: sketchybar.FontOptions{
				Font: settings.SketchybarSettings.IconFont,
				Kind: settings.SketchybarSettings.IconFontKind,
				Size: settings.SketchybarSettings.IconFontSize,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Color: sketchybar.ColorOptions{
				Color: settings.SketchybarSettings.LabelColor,
			},
			Font: sketchybar.FontOptions{
				Font: settings.SketchybarSettings.LabelFont,
				Kind: settings.SketchybarSettings.LabelFontKind,
				Size: settings.SketchybarSettings.LabelFontSize,
			},
			Padding: sketchybar.PaddingOptions{
				Right: 0,
				Left:  0,
			},
		},
		Background: sketchybar.BackgroundOptions{
			Drawing:      "on",
			Height:       settings.SketchybarSettings.ItemHeight,
			CornerRadius: 0,
			Color: sketchybar.ColorOptions{
				Color: settings.SketchybarSettings.ItemBackgroundColor,
			},
			Padding: sketchybar.PaddingOptions{
				Right: 0,
				Left:  0,
			},
		},
		Updates:     "off",
		ScrollTexts: "on",
	}

	batches = batch(batches, m(s("--default"), defaults.ToArgs()))

	return batches, nil
}
