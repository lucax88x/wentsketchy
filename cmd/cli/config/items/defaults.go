package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Defaults(batches [][]string) ([][]string, error) {
	defaults := sketchybar.ItemOptions{
		YOffset: pointer(0),
		Padding: sketchybar.PaddingOptions{
			Right: pointer(0),
			Left:  pointer(0),
		},
		Icon: sketchybar.ItemIconOptions{
			Color: sketchybar.ColorOptions{
				Color: settings.Sketchybar.IconColor,
			},
			Font: sketchybar.FontOptions{
				Font: settings.Sketchybar.IconFont,
				Kind: settings.Sketchybar.IconFontKind,
				Size: settings.Sketchybar.IconFontSize,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Color: sketchybar.ColorOptions{
				Color: settings.Sketchybar.LabelColor,
			},
			Font: sketchybar.FontOptions{
				Font: settings.Sketchybar.LabelFont,
				Kind: settings.Sketchybar.LabelFontKind,
				Size: settings.Sketchybar.LabelFontSize,
			},
			Padding: sketchybar.PaddingOptions{
				Right: pointer(0),
				Left:  pointer(0),
			},
		},
		Background: sketchybar.BackgroundOptions{
			Drawing:      "on",
			Height:       settings.Sketchybar.ItemHeight,
			CornerRadius: settings.Sketchybar.ItemRadius,
			Color: sketchybar.ColorOptions{
				Color: settings.Sketchybar.ItemBackgroundColor,
			},
			Border: sketchybar.BorderOptions{
				Color: settings.Sketchybar.ItemBorderColor,
				Width: pointer(2),
			},
			Padding: sketchybar.PaddingOptions{
				Right: pointer(0),
				Left:  pointer(0),
			},
		},
		Updates:     "off",
		ScrollTexts: "on",
	}

	batches = batch(batches, m(s("--default"), defaults.ToArgs()))

	return batches, nil
}
