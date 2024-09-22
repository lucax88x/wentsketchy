package items

import (
	"errors"
	"time"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/formatter"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type CalendarItem struct {
}

func NewCalendarItem() CalendarItem {
	return CalendarItem{}
}

const calendarItemName = "calendar"

func (i CalendarItem) Init(
	batches [][]string,
	fifoPath string,
) ([][]string, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("calendar: could not generate update event")
	}

	calendarItem := sketchybar.ItemOptions{
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
		Script:     updateEvent,
		// Click_script:            "$PLUGIN_DIR/zen.sh",
	}

	batches = batch(batches, s("--add", "item", calendarItemName, "right"))
	batches = batch(batches, m(s("--set", calendarItemName), calendarItem.ToArgs()))
	batches = batch(batches, s("--subscribe", calendarItemName, events.SystemWoke))

	return batches, nil
}

func (i CalendarItem) Update(
	batches [][]string,
	args *args.In,
) ([][]string, error) {
	if !isCalendar(args.Name) {
		return batches, nil
	}

	calendarItem := sketchybar.ItemOptions{
		Label: sketchybar.ItemLabelOptions{
			Value: formatter.HoursMinutes(time.Now()),
		},
	}

	batches = batch(batches, m(s("--set", calendarItemName), calendarItem.ToArgs()))

	return batches, nil
}

func isCalendar(name string) bool {
	return name == calendarItemName
}
