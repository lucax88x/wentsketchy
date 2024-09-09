package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

const font = "SF Pro"
const iconFont = "Hack Nerd Font"

func Init(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	err := initBar(ctx, di)

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

func initBar(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	bar := sketchybar.BarOptions{
		Position: "top",
		Height:   45,
		ColorOptions: sketchybar.ColorOptions{
			Color: ColorBarColor,
		},
		BorderOptions: sketchybar.BorderOptions{
			Width: 2,
			Color: ColorBarBorderColor,
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
