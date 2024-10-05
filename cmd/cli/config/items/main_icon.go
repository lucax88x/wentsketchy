package items

import (
	"context"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/icons"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type MainIconItem struct {
}

func NewMainIconItem() MainIconItem {
	return MainIconItem{}
}

const mainIconItemName = "main_icon"

func (i MainIconItem) Init(
	_ context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	mainIcon := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  settings.Sketchybar.ItemSpacing,
			Right: pointer(0),
		},
		Icon: sketchybar.ItemIconOptions{
			Value: icons.Apple,
			Color: sketchybar.ColorOptions{
				Color: colors.White,
			},
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	batches = batch(batches, s("--add", "item", mainIconItemName, position))
	batches = batch(batches, m(s("--set", mainIconItemName), mainIcon.ToArgs()))

	return batches, nil
}

func (i MainIconItem) Update(
	_ context.Context,
	batches Batches,
	_ sketchybar.Position,
	_ *args.In,
) (Batches, error) {
	return batches, nil
}

var _ WentsketchyItem = (*MainIconItem)(nil)
