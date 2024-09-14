package items

import (
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type FrontAppItem struct {
	AerospaceData *aerospace.Data
}

const frontAppItemName = "front_app"

func (i *FrontAppItem) Init(batches [][]string, fifoPath string) ([][]string, error) {
	frontApp := sketchybar.ItemOptions{
		Display: "active",
		// icon.background.drawing=on
		Icon: sketchybar.ItemIconOptions{
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  0,
			},
		},
		Script:      args.BuildEvent(fifoPath),
		ClickScript: "open -a 'Mission Control'",
	}

	batches = batch(batches, s("--add", "item", frontAppItemName, "right"))
	batches = batch(batches, m(s("--set", frontAppItemName), frontApp.ToArgs()))
	batches = batch(batches, s("--subscribe", frontAppItemName, string(events.FrontAppSwitched)))

	return batches, nil
}

func (i *FrontAppItem) Update(
	batches [][]string,
	args *args.Args,
) ([][]string, error) {
	if !isFrontApp(args.Name) {
		return batches, nil
	}

	if args.Event == string(events.FrontAppSwitched) {
		frontApp := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: args.Info,
			},

			Icon: sketchybar.ItemIconOptions{
				BackgroundOptions: sketchybar.BackgroundOptions{
					Drawing: true,
					ImageOptions: sketchybar.ImageOptions{
						Value:   fmt.Sprintf("app.%s", args.Info),
						Drawing: true,
						Scale:   "0.8",
					},
				},
			},
		}

		batches = batch(batches, m(s("--set", frontAppItemName), frontApp.ToArgs()))

		windows := i.AerospaceData.WindowsOfFocusedWorkspace("2")

		for _, window := range windows {
			if window.App != args.Info {
				continue
			}

			windowItemID := fmt.Sprintf("window.%d", window.ID)

			windowItem := sketchybar.ItemOptions{
				Icon: sketchybar.ItemIconOptions{
					Highlight: true,
				},
			}

			batches = batch(batches, m(s("--set", windowItemID), windowItem.ToArgs()))
		}
	}

	return batches, nil
}

func isFrontApp(name string) bool {
	return name == frontAppItemName
}
