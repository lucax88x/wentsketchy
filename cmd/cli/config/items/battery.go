package items

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/distatus/battery"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
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
	batches batches,
	fifoPath string,
) (batches, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("battery: could not generate update event")
	}

	batteryItem := sketchybar.ItemOptions{
		Padding: sketchybar.PaddingOptions{
			Left:  settings.SketchybarSettings.ItemSpacing,
			Right: settings.SketchybarSettings.ItemSpacing,
		},
		Background: sketchybar.BackgroundOptions{
			CornerRadius: settings.SketchybarSettings.ItemRadius,
		},
		Icon: sketchybar.ItemIconOptions{
			Value: settings.IconClock,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.SketchybarSettings.IconPadding,
				Right: settings.SketchybarSettings.IconPadding / 2,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			Padding: sketchybar.PaddingOptions{
				Left:  0,
				Right: settings.SketchybarSettings.IconPadding,
			},
		},
		UpdateFreq: 120,
		Updates:    "on",
		Script:     updateEvent,
	}

	batches = batch(batches, s("--add", "item", batteryItemName, "right"))
	batches = batch(batches, m(s("--set", batteryItemName), batteryItem.ToArgs()))
	batches = batch(batches, s("--subscribe", batteryItemName,
		events.PowerSourceChanged,
		events.SystemWoke,
	))

	return batches, nil
}

func (i BatteryItem) Update(
	batches batches,
	args *args.In,
) (batches, error) {
	if !isBattery(args.Name) {
		return batches, nil
	}

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

	return batches, nil
}

func isBattery(name string) bool {
	return name == batteryItemName
}

func getBatteryStatus(percentage float64) (string, string) {
	switch {
	case percentage >= 80 && percentage <= 100:
		return settings.IconBattery100, colors.Battery1
	case percentage >= 70 && percentage < 80:
		return settings.IconBattery75, colors.Battery2
	case percentage >= 40 && percentage < 70:
		return settings.IconBattery50, colors.Battery3
	case percentage >= 10 && percentage < 40:
		return settings.IconBattery25, colors.Battery4
	case percentage >= 0 && percentage < 10:
		return settings.IconBattery0, colors.Battery5
	default:
		return "", "" // Handle invalid percentages
	}
}

func getBatteryPercentage(battery *battery.Battery) float64 {
	return (battery.Current / battery.Full) * 100
}
