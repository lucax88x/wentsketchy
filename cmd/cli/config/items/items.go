package items

import (
	"context"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type WentsketchyItem interface {
	Init(
		ctx context.Context,
		position sketchybar.Position,
		batches Batches,
	) (Batches, error)
	Update(
		ctx context.Context,
		batches Batches,
		position sketchybar.Position,
		args *args.In,
	) (Batches, error)
}

type IndexedWentsketchyItems = map[string]WentsketchyItem

type WentsketchyItems struct {
	MainIcon  MainIconItem
	Calendar  CalendarItem
	FrontApp  FrontAppItem
	Aerospace AerospaceItem
	Battery   BatteryItem
	CPU       CPUItem
	Sensors   SensorsItem
}
