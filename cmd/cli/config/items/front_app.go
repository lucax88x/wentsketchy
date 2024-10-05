package items

import (
	"context"
	"errors"
	"fmt"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type FrontAppItem struct {
}

func NewFrontAppItem() FrontAppItem {
	return FrontAppItem{}
}

const frontAppItemName = "front_app"

func (i FrontAppItem) Init(
	_ context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	updateEvent, err := args.BuildEvent()

	if err != nil {
		return batches, errors.New("front_app: could not generate update event")
	}

	frontAppItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  settings.Sketchybar.ItemSpacing,
			Right: settings.Sketchybar.ItemSpacing,
		},
		Icon: sketchybar.ItemIconOptions{
			Background: sketchybar.BackgroundOptions{
				Drawing: "on",
				Image: sketchybar.ImageOptions{
					Drawing: "on",
					Padding: sketchybar.PaddingOptions{
						Left:  settings.Sketchybar.IconPadding,
						Right: pointer(*settings.Sketchybar.IconPadding / 2),
					},
				},
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Padding: sketchybar.PaddingOptions{
				Left:  pointer(0),
				Right: settings.Sketchybar.IconPadding,
			},
		},
		Updates:     "on",
		Script:      updateEvent,
		ClickScript: "open -a 'Mission Control'",
	}

	batches = batch(batches, s("--add", "item", frontAppItemName, position))
	batches = batch(batches, m(s("--set", frontAppItemName), frontAppItem.ToArgs()))
	batches = batch(batches, s("--subscribe", frontAppItemName, events.FrontAppSwitched))

	return batches, nil
}

func (i FrontAppItem) Update(
	ctx context.Context,
	batches Batches,
	_ sketchybar.Position,
	args *args.In,
) (Batches, error) {
	if !isFrontApp(args.Name) {
		return batches, nil
	}

	if args.Event == events.FrontAppSwitched {
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
	}

	return batches, nil
}

func isFrontApp(name string) bool {
	return name == frontAppItemName
}

var _ WentsketchyItem = (*FrontAppItem)(nil)
