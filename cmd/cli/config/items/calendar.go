package items

import (
	"time"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/formatter"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type CalendarItem struct {
}

const calendarItemName = "calendar"

func (i CalendarItem) Init(
	batches [][]string,
	fifoPath string,
) ([][]string, error) {
	calendar := sketchybar.ItemOptions{
		Icon: sketchybar.ItemIconOptions{
			Value: settings.IconClock,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
				Kind: "Regular",
				Size: "12.0",
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
		},
		Label: sketchybar.ItemLabelOptions{
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
		},
		Background: sketchybar.BackgroundOptions{
			BorderOptions: sketchybar.BorderOptions{
				Width: 2,
				Color: settings.ColorBackground1,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: settings.ColorBackground1,
			},
		},
		UpdateFreq: 30,
		Script:     args.BuildEvent(fifoPath),
		// Click_script:            "$PLUGIN_DIR/zen.sh",
	}

	batches = batch(batches, s("--add", "item", calendarItemName, "right"))
	batches = batch(batches, m(s("--set", calendarItemName), calendar.ToArgs()))
	batches = batch(batches, s("--subscribe", calendarItemName, "system_woke"))

	return batches, nil
}

func (i CalendarItem) Update(
	batches [][]string,
	args *args.Args,
) ([][]string, error) {
	if !isCalendar(args.Name) {
		return batches, nil
	}

	calendar := sketchybar.ItemOptions{
		Label: sketchybar.ItemLabelOptions{
			Value: formatter.HoursMinutes(time.Now()),
		},
	}

	batches = batch(batches, m(s("--set", calendarItemName), calendar.ToArgs()))

	return batches, nil
}

func isCalendar(name string) bool {
	return name == calendarItemName
}
