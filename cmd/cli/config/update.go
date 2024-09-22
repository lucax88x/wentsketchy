package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
)

func (cfg *Config) Update(
	ctx context.Context,
	args *args.In,
) error {
	var batches = make([][]string, 0)

	batches, err := cfg.items.Battery.Update(batches, args)

	if err != nil {
		return fmt.Errorf("update: battery %w", err)
	}

	batches, err = cfg.items.Calendar.Update(batches, args)

	if err != nil {
		return fmt.Errorf("update: calendar %w", err)
	}

	batches, err = cfg.items.FrontApp.Update(ctx, batches, args)

	if err != nil {
		return fmt.Errorf("update: calendar %w", err)
	}

	batches, err = cfg.items.Aerospace.Update(ctx, batches, args)

	if err != nil {
		return fmt.Errorf("update: aerospace %w", err)
	}

	err = cfg.sketchybar.Run(ctx, flatten(batches...))

	if err != nil {
		return fmt.Errorf("update: apply to sketchybar %w", err)
	}

	return nil
}
