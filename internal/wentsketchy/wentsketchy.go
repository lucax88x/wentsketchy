package wentsketchy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type Wentsketchy struct {
	Logger             *slog.Logger
	Clock              clock.Clock
	Sketchybar         sketchybar.API
	SketchybarSettings sketchybar.Settings
	Aerospace          aerospace.Tree
	aerospaceAPI       aerospace.API
}

func NewWentsketchy(
	ctx context.Context,
	logger *slog.Logger,
	clock clock.Clock,
) (*Wentsketchy, error) {
	di := &Wentsketchy{
		Logger: logger,
		Clock:  clock,
	}

	err := initialize(ctx, di)

	if err != nil {
		return nil, fmt.Errorf("init: could not initialize wentsketchy. %w", err)
	}

	return di, nil
}

func initialize(_ context.Context, di *Wentsketchy) error {
	di.aerospaceAPI = aerospace.NewAPI(di.Logger)
	di.Aerospace = aerospace.NewTree(di.Logger, di.aerospaceAPI)
	di.Sketchybar = sketchybar.NewAPI(di.Logger)

	di.SketchybarSettings = sketchybar.Settings{
		LabelColor:    ColorWhite,
		LabelFont:     FontLabel,
		LabelFontKind: "Semibold",
		LabelFontSize: "14.0",
		IconColor:     ColorWhite,
		IconFont:      FontIcon,
		IconFontKind:  "Regular",
		IconFontSize:  "14.0",
		IconStripFont: FontAppIcon,
	}

	return nil
}
