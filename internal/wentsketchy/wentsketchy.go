package wentsketchy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/clock"
	"github.com/lucax88x/wentsketchy/internal/command"
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
	command              *command.Command
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
	cfg, err := config.ReadYaml()

	if err != nil {
		//nolint:errorlint // no wrap
		return fmt.Errorf("wentsketchy: could not initialize cfg from yaml. %v", err)
	}

	di.command = command.NewCommand(di.Logger)
	di.aerospaceAPI = aerospace.NewAPI(di.Logger, di.command)
	di.aerospaceTreeBuilder = aerospace.NewTreeBuilder(di.Logger, di.aerospaceAPI)
	di.Aerospace = aerospace.New(di.Logger, di.aerospaceAPI, di.aerospaceTreeBuilder)

	di.Sketchybar = sketchybar.NewAPI(di.Logger, di.command)

	mainIcon := items.NewMainIconItem()
	calendar := items.NewCalendarItem()
	frontApp := items.NewFrontAppItem()
	aerospace := items.NewAerospaceItem(di.Logger, di.Aerospace, di.Sketchybar)
	battery := items.NewBatteryItem(di.Logger)
	cpu := items.NewCPUItem(di.Logger, di.command)
	sensors := items.NewSensorsItem(di.Logger, di.command)

	di.Config = config.NewConfig(
		cfg,
		di.Logger,
		di.Sketchybar,
		map[string]items.WentsketchyItem{
			"main_icon": mainIcon,
			"calendar":  calendar,
			"front_app": frontApp,
			"aerospace": aerospace,
			"battery":   battery,
			"cpu":       cpu,
			"sensors":   sensors,
		},
		items.WentsketchyItems{
			MainIcon:  mainIcon,
			Calendar:  calendar,
			FrontApp:  frontApp,
			Aerospace: aerospace,
			Battery:   battery,
			CPU:       cpu,
			Sensors:   sensors,
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
