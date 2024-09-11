package config

import (
	"context"
	"fmt"

	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func initLeft(ctx context.Context, di *wentsketchy.Wentsketchy) error {
	batches, err := initAerospace(ctx, di)

	if err != nil {
		return fmt.Errorf("init aerospace %w", err)
	}

	err = di.Sketchybar.Run(ctx, batches)

	if err != nil {
		return fmt.Errorf("apply to sketchybar %w", err)
	}

	return nil
}
