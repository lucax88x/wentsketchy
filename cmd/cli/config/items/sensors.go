package items

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/icons"
	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

const statsApp = "/Applications/Stats.app/Contents/Resources/smc"

type SensorsItem struct {
	logger  *slog.Logger
	command *command.Command
}

func NewSensorsItem(logger *slog.Logger, command *command.Command) SensorsItem {
	return SensorsItem{
		logger,
		command,
	}
}

const sensorsBracketName = "sensors.bracket"
const sensorsItemIconName = "sensors.icon"
const sensorsItemFansName = "sensors.fans"
const sensorsItemTemperaturesName = "sensors.temperatures"
const sensorsItemSpacerName = "sensors.spacer"

func (i SensorsItem) Init(
	_ context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	updateEvent, err := args.BuildEvent()

	if err != nil {
		return batches, errors.New("sensors: could not generate update event")
	}

	sensorsIconItem := sketchybar.ItemOptions{
		Display: "active",
		Icon: sketchybar.ItemIconOptions{
			Value: icons.ThermoMedium,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.Sketchybar.IconPadding,
				Right: pointer(*settings.Sketchybar.IconPadding / 2),
			},
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	sensorsFansItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  pointer(0),
			Right: settings.Sketchybar.ItemSpacing,
		},
		Label: sketchybar.ItemLabelOptions{
			Value: "",
			Font: sketchybar.FontOptions{
				Size: "8.0",
			},
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
		YOffset:    pointer(-6),
		Width:      pointer(0),
		UpdateFreq: pointer(4),
		Updates:    "on",
		Script:     updateEvent,
	}
	sensorsTemperaturesItem := sketchybar.ItemOptions{
		Display: "active",
		Padding: sketchybar.PaddingOptions{
			Left:  pointer(0),
			Right: settings.Sketchybar.ItemSpacing,
		},
		Label: sketchybar.ItemLabelOptions{
			Value: "",
			Font: sketchybar.FontOptions{
				Size: "8.0",
			},
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: "off",
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
		YOffset: pointer(4),
		// Width:   pointer(0),
	}
	sensorsBracketItem := sketchybar.BracketOptions{
		Background: sketchybar.BackgroundOptions{
			Drawing: "on",
		},
	}
	sensorsSpacerItem := sketchybar.ItemOptions{
		Display: "active",
		Label: sketchybar.ItemLabelOptions{
			Value: "",
		},
		Padding: sketchybar.PaddingOptions{
			Right: settings.Sketchybar.ItemSpacing,
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	batches = batch(batches, s("--add", "item", sensorsItemSpacerName, position))
	batches = batch(batches, m(s("--set", sensorsItemSpacerName), sensorsSpacerItem.ToArgs()))

	batches = batch(batches, s("--add", "item", sensorsItemFansName, position))
	batches = batch(batches, m(s("--set", sensorsItemFansName), sensorsFansItem.ToArgs()))

	batches = batch(batches, s("--add", "item", sensorsItemTemperaturesName, position))
	batches = batch(batches, m(s("--set", sensorsItemTemperaturesName), sensorsTemperaturesItem.ToArgs()))

	batches = batch(batches, s("--add", "item", sensorsItemIconName, position))
	batches = batch(batches, m(s("--set", sensorsItemIconName), sensorsIconItem.ToArgs()))

	batches = batch(batches, s(
		"--add",
		"bracket",
		sensorsBracketName,
		sensorsItemIconName,
		sensorsItemFansName,
		sensorsItemTemperaturesName,
	))
	batches = batch(batches, m(s("--set", sensorsBracketName), sensorsBracketItem.ToArgs()))

	return batches, nil
}

func (i SensorsItem) Update(
	ctx context.Context,
	batches Batches,
	_ sketchybar.Position,
	args *args.In,
) (Batches, error) {
	if !isFAN(args.Name) {
		return batches, nil
	}

	if args.Event == events.Routine || args.Event == events.Forced {
		fanSpeeds, err := i.getFanSpeeds(ctx)

		if err != nil {
			return batches, err
		}

		temperatures, err := i.getTemperatures(ctx)

		if err != nil {
			return batches, err
		}

		fanSpeed := float32(-1)
		if len(fanSpeeds) > 0 {
			fanSpeed = fanSpeeds[0]
		}

		actualFanSpeed := "Fans Off"

		if fanSpeed > -1 {
			actualFanSpeed = fmt.Sprintf("%.0f RPM", fanSpeed)
		}

		sensorsFanItem := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: actualFanSpeed,
			},
		}
		sensorsTemperaturesItem := sketchybar.ItemOptions{
			Label: sketchybar.ItemLabelOptions{
				Value: fmt.Sprintf("%.0f°C / %.0f°C", temperatures.highest, temperatures.averageCPUs),
			},
		}
		batches = batch(batches, m(s("--set", sensorsItemFansName), sensorsFanItem.ToArgs()))
		batches = batch(batches, m(s("--set", sensorsItemTemperaturesName), sensorsTemperaturesItem.ToArgs()))
	}

	return batches, nil
}

func isFAN(name string) bool {
	return name == sensorsItemFansName
}

func (i SensorsItem) getFanSpeeds(ctx context.Context) ([]float32, error) {
	out, err := i.command.Run(ctx, statsApp, "fans")

	if err != nil {
		return make([]float32, 0), fmt.Errorf("sensors: could not get fan speed. %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(out))
	results := make([]float32, 0)
	for scanner.Scan() {
		line := scanner.Text()

		speedFromLine, cut := strings.CutPrefix(line, "Actual speed: ")
		if cut {
			conv, err := strconv.ParseFloat(speedFromLine, 32)

			if err != nil {
				//nolint:errorlint // no wrap
				return make([]float32, 0), fmt.Errorf("sensors: could not parse fan speed from line %s. %v", line, err)
			}

			results = append(results, float32(conv))
		}
	}

	return results, nil
}

type temperatures struct {
	highest     float32
	averageCPUs float32
}

func (i SensorsItem) getTemperatures(ctx context.Context) (temperatures, error) {
	var results temperatures
	out, err := i.command.Run(ctx, statsApp, "list", "-t")

	if err != nil {
		return results, fmt.Errorf("sensors: could not get fan speed. %w", err)
	}

	var cpuTemps []float32
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "[INFO]") {
			continue
		}

		if len(line) == 0 {
			continue
		}

		temp, err := parseTemperature(line)

		if err != nil {
			continue
		}

		if temp <= 0 {
			continue
		}

		if temp > results.highest {
			results.highest = temp
		}

		if strings.HasPrefix(line, "[TC") {
			cpuTemps = append(cpuTemps, temp)
		}
	}

	if len(cpuTemps) > 0 {
		sum := float32(0)
		for _, temp := range cpuTemps {
			sum += temp
		}
		results.averageCPUs = sum / float32(len(cpuTemps))
	}

	return results, nil
}

func parseTemperature(line string) (float32, error) {
	parts := strings.Fields(line)

	if len(parts) < 2 {
		return 0, fmt.Errorf("sensors: invalid temperature line format with %s", line)
	}

	part := parts[len(parts)-1]

	temp, err := strconv.ParseFloat(part, 32)
	if err != nil {
		//nolint:errorlint // no wrap
		return 0, fmt.Errorf("sensors failed to parse temperature from %s: %v", part, err)
	}
	return float32(temp), nil
}

var _ WentsketchyItem = (*SensorsItem)(nil)
