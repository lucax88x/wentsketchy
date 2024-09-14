package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type Config struct {
	logger     *slog.Logger
	sketchybar sketchybar.API
	items      items.WentsketchyItems
	fifoPath   string
}

func NewConfig(
	logger *slog.Logger,
	sketchybar sketchybar.API,
	items items.WentsketchyItems,
) *Config {
	return &Config{
		logger,
		sketchybar,
		items,
		"",
	}
}

func (cfg *Config) SetFifoPath(fifoPath string) {
	cfg.fifoPath = fifoPath
}

func (cfg *Config) Init(ctx context.Context) error {
	var batches = make([][]string, 0)

	batches, err := items.Defaults(batches)

	if err != nil {
		return fmt.Errorf("config: defaults %w", err)
	}

	batches, err = items.Bar(batches)

	if err != nil {
		return fmt.Errorf("config: bar %w", err)
	}

	batches, err = cfg.left(batches)

	if err != nil {
		return fmt.Errorf("config: left %w", err)
	}

	batches, err = cfg.right(batches)

	if err != nil {
		return fmt.Errorf("config: right %w", err)
	}

	err = cfg.sketchybar.Run(ctx, flatten(batches...))

	if err != nil {
		return fmt.Errorf("config: apply to sketchybar %w", err)
	}

	err = cfg.sketchybar.Run(ctx, []string{"--update"})

	if err != nil {
		return fmt.Errorf("config: update sketchybar %w", err)
	}

	return nil
}
