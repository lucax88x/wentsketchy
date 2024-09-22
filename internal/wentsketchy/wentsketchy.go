package wentsketchy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/lucax88x/wentsketchy/internal/fifo"
	"github.com/lucax88x/wentsketchy/internal/server"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type Wentsketchy struct {
	Logger               *slog.Logger
	Clock                clock.Clock
	Config               *config.Config
	Fifo                 *fifo.Reader
	Server               *server.FifoServer
	Sketchybar           sketchybar.API
	Aerospace            aerospace.Aerospace
	aerospaceTreeBuilder aerospace.TreeBuilder
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
	di.aerospaceTreeBuilder = aerospace.NewTreeBuilder(di.Logger, di.aerospaceAPI)
	di.Aerospace = aerospace.New(di.Logger, di.aerospaceAPI, di.aerospaceTreeBuilder)

	di.Sketchybar = sketchybar.NewAPI(di.Logger)
	di.Config = config.NewConfig(
		di.Logger,
		di.Sketchybar,
		items.WentsketchyItems{
			MainIcon: items.NewMainIconItem(),
			Aerospace: items.NewAerospaceItem(
				di.Aerospace,
				di.Sketchybar,
			),
			Calendar: items.NewCalendarItem(),
			FrontApp: items.NewFrontAppItem(di.Aerospace),
			Battery:  items.NewBatteryItem(di.Logger),
		},
	)

	di.Fifo = fifo.NewFifoReader(di.Logger)
	di.Server = server.NewFifoServer(
		di.Logger,
		di.Config,
		di.Fifo,
		di.Aerospace,
	)

	return nil
}
