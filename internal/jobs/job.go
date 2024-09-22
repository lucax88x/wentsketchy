package jobs

import (
	"context"
	"time"

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
			di.Aerospace.RefreshTree()
		case <-ctx.Done():
			di.Logger.InfoContext(ctx, "jobs: cancel")
			return
		}
	}
}
