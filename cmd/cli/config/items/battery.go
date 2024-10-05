package items

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/distatus/battery"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/icons"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type BatteryItem struct {
	logger *slog.Logger
}

func NewBatteryItem(logger *slog.Logger) BatteryItem {
	return BatteryItem{logger}
}

const batteryItemName = "battery"

func (i BatteryItem) Init(
	_ context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	updateEvent, err := args.BuildEvent()

	if err != nil {
		return batches, errors.New("battery: could not generate update event")
	}

	batteryItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  settings.Sketchybar.ItemSpacing,
			Right: settings.Sketchybar.ItemSpacing,
		},
		Icon: sketchybar.ItemIconOptions{
			Value: icons.Battery100,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.Sketchybar.IconPadding,
				Right: pointer(*settings.Sketchybar.IconPadding / 2),
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Padding: sketchybar.PaddingOptions{
				Left:  pointer(0),
				Right: settings.Sketchybar.IconPadding,
			},
		},
		UpdateFreq: pointer(120),
		Updates:    "on",
		Script:     updateEvent,
	}

	batches = batch(batches, s("--add", "item", batteryItemName, position))
	batches = batch(batches, m(s("--set", batteryItemName), batteryItem.ToArgs()))
	batches = batch(batches, s("--subscribe", batteryItemName,
		events.PowerSourceChanged,
		events.SystemWoke,
	))

	return batches, nil
}

func (i BatteryItem) Update(
	_ context.Context,
	batches Batches,
	_ sketchybar.Position,
	args *args.In,
) (Batches, error) {
	if !isBattery(args.Name) {
		return batches, nil
	}

	if args.Event == events.Routine || args.Event == events.Forced {
		batteries, err := battery.GetAll()

		if err != nil {
			return batches, fmt.Errorf("battery: could not get battery info. %w", err)
		}

		if len(batteries) == 0 {
			return batches, errors.New("battery: has no battery")
		}

		if len(batteries) > 1 {
			i.logger.Warn(
				"does not support multiple batteries",
				slog.Int("batteries", len(batteries)),
			)
		}

		battery := batteries[0]

		percentage := getBatteryPercentage(battery)
		icon, color := getBatteryStatus(percentage)

		batteryItem := sketchybar.ItemOptions{
			Icon: sketchybar.ItemIconOptions{
				Value: icon,
				Color: sketchybar.ColorOptions{
					Color: color,
				},
			},
			Label: sketchybar.ItemLabelOptions{
				Value: fmt.Sprintf("%.0f%%", percentage),
			},
		}

		batches = batch(batches, m(s("--set", batteryItemName), batteryItem.ToArgs()))
	}

	return batches, nil
}

func isBattery(name string) bool {
	return name == batteryItemName
}

func getBatteryStatus(percentage float64) (string, string) {
	switch {
	case percentage >= 80 && percentage <= 100:
		return icons.Battery100, colors.Battery1
	case percentage >= 70 && percentage < 80:
		return icons.Battery75, colors.Battery2
	case percentage >= 40 && percentage < 70:
		return icons.Battery50, colors.Battery3
	case percentage >= 10 && percentage < 40:
		return icons.Battery25, colors.Battery4
	case percentage >= 0 && percentage < 10:
		return icons.Battery0, colors.Battery5
	default:
		return "", "" // Handle invalid percentages
	}
}

func getBatteryPercentage(battery *battery.Battery) float64 {
	return (battery.Current / battery.Full) * 100
}

var _ WentsketchyItem = (*BatteryItem)(nil)
