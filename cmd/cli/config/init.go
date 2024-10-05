package config

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type Config struct {
	Cfg          *Cfg
	logger       *slog.Logger
	sketchybar   sketchybar.API
	IndexedItems items.IndexedWentsketchyItems
	Items        items.WentsketchyItems
}

func NewConfig(
	cfg *Cfg,
	logger *slog.Logger,
	sketchybar sketchybar.API,
	indexedItems items.IndexedWentsketchyItems,
	items items.WentsketchyItems,
) *Config {
	return &Config{
		cfg,
		logger,
		sketchybar,
		indexedItems,
		items,
	}
}

func (cfg *Config) Init(ctx context.Context) error {
	var batches = make(items.Batches, 0)

	batches, err := items.Defaults(batches)

	if err != nil {
		return fmt.Errorf("config: defaults %w", err)
	}

	batches, err = items.Bar(batches)

	if err != nil {
		return fmt.Errorf("config: bar %w", err)
	}

	batches, err = cfg.initList(ctx, batches, sketchybar.PositionLeft, cfg.Cfg.Left)

	if err != nil {
		return fmt.Errorf("config: left %w", err)
	}

	batches, err = cfg.initList(ctx, batches, sketchybar.PositionLeftNotch, cfg.Cfg.LeftNotch)

	if err != nil {
		return fmt.Errorf("config: left notch %w", err)
	}

	batches, err = cfg.initList(ctx, batches, sketchybar.PositionCenter, cfg.Cfg.Center)

	if err != nil {
		return fmt.Errorf("config: center %w", err)
	}

	batches, err = cfg.initList(ctx, batches, sketchybar.PositionRight, reverse(cfg.Cfg.Right))

	if err != nil {
		return fmt.Errorf("config: right %w", err)
	}

	batches, err = cfg.initList(ctx, batches, sketchybar.PositionRightNotch, reverse(cfg.Cfg.RightNotch))

	if err != nil {
		return fmt.Errorf("config: right notch %w", err)
	}

	err = cfg.sketchybar.Run(ctx, items.Flatten(batches...))

	if err != nil {
		return fmt.Errorf("config: apply to sketchybar %w", err)
	}

	batches = make(items.Batches, 0)
	batches, err = items.ShowBar(batches)

	if err != nil {
		return fmt.Errorf("config: appear bar %w", err)
	}

	err = cfg.sketchybar.Run(ctx, items.Flatten(batches...))

	if err != nil {
		return fmt.Errorf("config: apply to sketchybar %w", err)
	}

	err = cfg.sketchybar.Run(ctx, []string{"--update"})

	if err != nil {
		return fmt.Errorf("config: update sketchybar %w", err)
	}

	return nil
}

func (cfg *Config) initList(
	ctx context.Context,
	batches items.Batches,
	position sketchybar.Position,
	list []string,
) (items.Batches, error) {
	var err error
	for _, itemName := range list {
		item, found := cfg.IndexedItems[itemName]

		if found {
			batches, err = item.Init(ctx, position, batches)

			if err != nil {
				return batches, fmt.Errorf("init: error while init %s. %w", itemName, err)
			}
		} else {
			return batches, fmt.Errorf("init: did not find %s", itemName)
		}
	}

	return batches, nil
}

func reverse(items []string) []string {
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}
	return items
}
