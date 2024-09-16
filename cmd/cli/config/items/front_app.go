package items

import (
	"errors"
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type FrontAppItem struct {
	AerospaceData *aerospace.Data
}

func NewFrontAppItem(data *aerospace.Data) FrontAppItem {
	return FrontAppItem{data}
}

const frontAppItemName = "front_app"

func (i *FrontAppItem) Init(batches [][]string, fifoPath string) ([][]string, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("front_app: could not generate update event")
	}

	frontAppItem := sketchybar.ItemOptions{
		Display: "active",
		Background: sketchybar.BackgroundOptions{
			Drawing: true,
		},
		Icon: sketchybar.ItemIconOptions{
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  0,
			},
		},
		Updates:     "on",
		Script:      updateEvent,
		ClickScript: "open -a 'Mission Control'",
	}

	batches = batch(batches, s("--add", "item", frontAppItemName, "right"))
	batches = batch(batches, m(s("--set", frontAppItemName), frontAppItem.ToArgs()))
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
		frontAppItem := sketchybar.ItemOptions{
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

		batches = batch(batches, m(s("--set", frontAppItemName), frontAppItem.ToArgs()))

		windowsOfPrevWorkspace := i.AerospaceData.WindowsOfWorkspace(i.AerospaceData.PrevWorkspaceID)

		for _, window := range windowsOfPrevWorkspace {
			windowItemID := fmt.Sprintf("window.%d", window.ID)

			windowItem := sketchybar.ItemOptions{
				Icon: sketchybar.ItemIconOptions{
					Highlight: false,
				},
			}

			batches = batch(batches, m(s("--set", windowItemID), windowItem.ToArgs()))
		}

		windowsOfFocusedWorkspace := i.AerospaceData.WindowsOfWorkspace(i.AerospaceData.FocusedWorkspaceID)

		for _, window := range windowsOfFocusedWorkspace {
			windowItemID := fmt.Sprintf("window.%d", window.ID)

			windowItem := sketchybar.ItemOptions{
				Icon: sketchybar.ItemIconOptions{
					Highlight: window.App == args.Info,
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
