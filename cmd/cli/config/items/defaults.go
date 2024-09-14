package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func Defaults(batches [][]string) ([][]string, error) {
	defaults := sketchybar.ItemOptions{
		//   updates=when_shown
		//   scroll_texts=on
		//   padding_right="$PADDINGS"
		//   padding_left="$PADDINGS"
		// PaddingOptions: sketchybar.PaddingOptions{
		// 	Right: 4,
		// 	Left:  4,
		// },
		Icon: sketchybar.ItemIconOptions{
			ColorOptions: sketchybar.ColorOptions{
				Color: settings.SketchybarSettings.IconColor,
			},
			Font: sketchybar.FontOptions{
				Font: settings.SketchybarSettings.IconFont,
				Kind: settings.SketchybarSettings.IconFontKind,
				Size: settings.SketchybarSettings.IconFontSize,
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 4,
				Left:  4,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			ColorOptions: sketchybar.ColorOptions{
				Color: settings.SketchybarSettings.LabelColor,
			},
			//   label.shadow.drawing=on
			//   label.shadow.distance=2
			//   label.shadow.color=0xff000000
			Font: sketchybar.FontOptions{
				Font: settings.SketchybarSettings.LabelFont,
				Kind: settings.SketchybarSettings.LabelFontKind,
				Size: settings.SketchybarSettings.IconFontSize,
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 4,
				Left:  4,
			},
		},
		Background: sketchybar.BackgroundOptions{
			//   background.corner_radius=4
			//   background.height=26
			BorderOptions: sketchybar.BorderOptions{
				Width: 2,
			},
		},
	}

	batches = batch(batches, m(s("--default"), defaults.ToArgs()))

	return batches, nil
}
