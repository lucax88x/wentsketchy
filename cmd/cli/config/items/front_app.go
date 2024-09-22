package items

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
	"github.com/lucax88x/wentsketchy/internal/utils"
)

type FrontAppItem struct {
	Aerospace aerospace.Aerospace
}

func NewFrontAppItem(data aerospace.Aerospace) FrontAppItem {
	return FrontAppItem{data}
}

const frontAppItemName = "front_app"

func (i *FrontAppItem) Init(batches batches, fifoPath string) (batches, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("front_app: could not generate update event")
	}

	frontAppItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  settings.SketchybarSettings.ItemSpacing,
			Right: settings.SketchybarSettings.ItemSpacing,
		},
		Background: sketchybar.BackgroundOptions{
			CornerRadius: settings.SketchybarSettings.ItemRadius,
		},
		Icon: sketchybar.ItemIconOptions{
			Background: sketchybar.BackgroundOptions{
				Drawing: "on",
				Image: sketchybar.ImageOptions{
					Drawing: "on",
					Padding: sketchybar.PaddingOptions{
						Left:  settings.SketchybarSettings.IconPadding,
						Right: settings.SketchybarSettings.IconPadding / 2,
					},
				},
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Padding: sketchybar.PaddingOptions{
				Left:  0,
				Right: settings.SketchybarSettings.IconPadding,
			},
		},
		Updates:     "on",
		Script:      updateEvent,
		ClickScript: "open -a 'Mission Control'",
	}

	batches = batch(batches, s("--add", "item", frontAppItemName, "right"))
	batches = batch(batches, m(s("--set", frontAppItemName), frontAppItem.ToArgs()))
	batches = batch(batches, s("--subscribe", frontAppItemName, events.FrontAppSwitched))

	return batches, nil
}

func (i *FrontAppItem) Update(
	ctx context.Context,
	batches batches,
	args *args.In,
) (batches, error) {
	if !isFrontApp(args.Name) {
		return batches, nil
	}

	if args.Event == events.FrontAppSwitched {
		i.Aerospace.SetFocusedApp(args.Info)

		frontAppItem := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: args.Info,
			},
			Icon: sketchybar.ItemIconOptions{
				Background: sketchybar.BackgroundOptions{
					Image: sketchybar.ImageOptions{
						Value: fmt.Sprintf("app.%s", args.Info),
						Scale: "0.8",
					},
				},
			},
		}

		batches = batch(batches, m(s("--set", frontAppItemName), frontAppItem.ToArgs()))

		tree := i.Aerospace.GetTree()

		batches = i.removeAllHighlights(batches, tree)
		batches = i.highlightWindows(ctx, batches, args.Info)
	}

	return batches, nil
}

func (i *FrontAppItem) removeAllHighlights(batches batches, tree *aerospace.Tree) batches {
	for _, window := range tree.IndexedWindows {
		windowItemID := fmt.Sprintf("window.%d", window.ID)

		windowItem := sketchybar.ItemOptions{
			Background: sketchybar.BackgroundOptions{
				Color: sketchybar.ColorOptions{
					Color: settings.SketchybarSettings.ItemBackgroundColor,
				},
			},
		}

		batches = batch(batches, m(s("--set", windowItemID), windowItem.ToArgs()))
	}

	return batches
}

func (i *FrontAppItem) highlightWindows(ctx context.Context, batches batches, app string) batches {
	windowsOfFocusedWorkspace := i.Aerospace.WindowsOfWorkspace(i.Aerospace.GetFocusedWorkspaceID(ctx))

	for _, window := range windowsOfFocusedWorkspace {
		windowItemID := fmt.Sprintf("window.%d", window.ID)

		backgroundColor := settings.SketchybarSettings.ItemBackgroundColor
		if utils.Equals(window.App, app) {
			backgroundColor = settings.SketchybarSettings.AerospaceItemFocusedBackgroundColor
		}

		windowItem := sketchybar.ItemOptions{
			Background: sketchybar.BackgroundOptions{
				Color: sketchybar.ColorOptions{
					Color: backgroundColor,
				},
			},
		}

		batches = batch(batches, m(s("--set", windowItemID), windowItem.ToArgs()))
	}

	return batches
}

func isFrontApp(name string) bool {
	return name == frontAppItemName
}
