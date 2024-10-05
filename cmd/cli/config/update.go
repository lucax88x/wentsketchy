package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

func (cfg *Config) Update(
	ctx context.Context,
	args *args.In,
) error {
	var batches = make(items.Batches, 0)

	batches, err := cfg.updateList(ctx, batches, sketchybar.PositionLeft, args, cfg.Cfg.Left)

	if err != nil {
		return fmt.Errorf("update: left %w", err)
	}

	batches, err = cfg.updateList(ctx, batches, sketchybar.PositionLeftNotch, args, cfg.Cfg.LeftNotch)

	if err != nil {
		return fmt.Errorf("update: left notch %w", err)
	}

	batches, err = cfg.updateList(ctx, batches, sketchybar.PositionCenter, args, cfg.Cfg.Center)

	if err != nil {
		return fmt.Errorf("update: center %w", err)
	}

	batches, err = cfg.updateList(ctx, batches, sketchybar.PositionRight, args, reverse(cfg.Cfg.Right))

	if err != nil {
		return fmt.Errorf("update: right %w", err)
	}

	batches, err = cfg.updateList(ctx, batches, sketchybar.PositionRightNotch, args, reverse(cfg.Cfg.RightNotch))

	if err != nil {
		return fmt.Errorf("update: right notch %w", err)
	}

	err = cfg.sketchybar.Run(ctx, items.Flatten(batches...))

	if err != nil {
		return fmt.Errorf("update: apply to sketchybar %w", err)
	}

	return nil
}

func (cfg *Config) updateList(
	ctx context.Context,
	batches items.Batches,
	position sketchybar.Position,
	args *args.In,
	list []string,
) (items.Batches, error) {
	var err error
	for _, itemName := range list {
		item, found := cfg.IndexedItems[itemName]

		if found {
			batches, err = item.Update(ctx, batches, position, args)

			if err != nil {
				return batches, fmt.Errorf("init: error while init %s. %w", itemName, err)
			}
		} else {
			return batches, fmt.Errorf("init: did not find %s", itemName)
		}
	}
	return batches, nil
}
