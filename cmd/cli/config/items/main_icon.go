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

const mainIconItemName = "main_icon"

func (i MainIconItem) Init(batches [][]string) ([][]string, error) {
	mainIcon := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  settings.SketchybarSettings.ItemSpacing,
			Right: 0,
		},
		Icon: sketchybar.ItemIconOptions{
			Value: settings.IconApple,
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	batches = batch(batches, s("--add", "item", mainIconItemName, "left"))
	batches = batch(batches, m(s("--set", mainIconItemName), mainIcon.ToArgs()))

	return batches, nil
}
