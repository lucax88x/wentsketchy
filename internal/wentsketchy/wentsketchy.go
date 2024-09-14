package wentsketchy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/lucax88x/wentsketchy/internal/server"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type Wentsketchy struct {
	Logger               *slog.Logger
	Clock                clock.Clock
	Config               *config.Config
	Server               server.FifoServer
	Sketchybar           sketchybar.API
	RefreshAerospaceData func()
	AerospaceData        *aerospace.Data
	Aerospace            aerospace.Tree
	aerospaceAPI         aerospace.API
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

func initialize(ctx context.Context, di *Wentsketchy) error {
	di.aerospaceAPI = aerospace.NewAPI(di.Logger)
	di.Aerospace = aerospace.NewTree(di.Logger, di.aerospaceAPI)
	di.AerospaceData = &aerospace.Data{}

	di.RefreshAerospaceData = func() {
		tree, err := di.Aerospace.Build()

		if err != nil {
			di.Logger.ErrorContext(ctx, "aerospace: could not refresh tree")
			return
		}

		di.AerospaceData.Tree = tree
	}

	di.Sketchybar = sketchybar.NewAPI(di.Logger)
	di.Config = config.NewConfig(
		di.Logger,
		di.Sketchybar,
		items.WentsketchyItems{
			Aerospace: items.AerospaceItem{AerospaceData: di.AerospaceData},
			Calendar:  items.CalendarItem{},
			FrontApp:  items.FrontAppItem{AerospaceData: di.AerospaceData},
		},
	)

	di.Server = server.NewFifoServer(di.Logger, di.Config)

	return nil
}
