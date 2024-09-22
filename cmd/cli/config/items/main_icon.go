package items

import (
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type MainIconItem struct {
}

func NewMainIconItem() MainIconItem {
	return MainIconItem{}
}

const mainIconItemName = "mainIcon"

func (i MainIconItem) Init(batches [][]string) ([][]string, error) {
	mainIcon := sketchybar.ItemOptions{
		// Display: "active",
		Icon: sketchybar.ItemIconOptions{
			Value: settings.IconApple,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
				Kind: "Regular",
				Size: "16.0",
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 8,
			},
		},
	}

	batches = batch(batches, s("--add", "item", mainIconItemName, "left"))
	batches = batch(batches, m(s("--set", mainIconItemName), mainIcon.ToArgs()))

	return batches, nil
}
