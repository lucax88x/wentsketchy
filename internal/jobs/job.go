package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/items"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func RefreshAerospaceData(
	ctx context.Context,
	di *wentsketchy.Wentsketchy,
	d time.Duration,
) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			di.Logger.InfoContext(ctx, "jobs: refreshing aerospace tree")

			var batches = make(items.Batches, 0)

			batches, err := di.Config.Items.Aerospace.CheckTree(ctx, batches)

			if err != nil {
				di.Logger.ErrorContext(ctx, "jobs: error checking tree.", slog.Any("err", err))
				return
			}

			err = di.Sketchybar.Run(ctx, items.Flatten(batches...))

			if err != nil {
				di.Logger.ErrorContext(ctx, "jobs: error while running sketchybar.", slog.Any("err", err))
				return
			}

			di.Logger.InfoContext(ctx, "jobs: checked aerospace tree")
		case <-ctx.Done():
			di.Logger.InfoContext(ctx, "jobs: cancel")
			return
		}
	}
}
