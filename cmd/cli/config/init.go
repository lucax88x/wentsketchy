package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func Init(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	err := defaults(ctx, di)

	if err != nil {
		return fmt.Errorf("defaults %w", err)
	}

	err = initBar(ctx, di)

	if err != nil {
		return fmt.Errorf("init bar %w", err)
	}

	err = initLeft(ctx, di)

	if err != nil {
		return fmt.Errorf("init left %w", err)
	}

	err = initRight(ctx, di)

	if err != nil {
		return fmt.Errorf("init right %w", err)
	}

	return nil
}

func defaults(ctx context.Context, di *wentsketchy.Wentsketchy) error {
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
				Color: di.SketchybarSettings.IconColor,
			},
			Font: sketchybar.FontOptions{
				Font: di.SketchybarSettings.IconFont,
				Kind: di.SketchybarSettings.IconFontKind,
				Size: di.SketchybarSettings.IconFontSize,
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 4,
				Left:  4,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			ColorOptions: sketchybar.ColorOptions{
				Color: di.SketchybarSettings.LabelColor,
			},
			//   label.shadow.drawing=on
			//   label.shadow.distance=2
			//   label.shadow.color=0xff000000
			Font: sketchybar.FontOptions{
				Font: di.SketchybarSettings.LabelFont,
				Kind: di.SketchybarSettings.LabelFontKind,
				Size: di.SketchybarSettings.IconFontSize,
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

	return di.Sketchybar.Run(ctx, m(s("--default"), defaults.ToArgs()))
}

func initBar(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   45,
		// TODO: move to settings
		ColorOptions: sketchybar.ColorOptions{
			Color: wentsketchy.ColorBarColor,
		},
		BorderOptions: sketchybar.BorderOptions{
			Width: 2,
			Color: wentsketchy.ColorBarBorderColor,
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

	return di.Sketchybar.Run(ctx, m(s("--bar"), bar.ToArgs()))
}
