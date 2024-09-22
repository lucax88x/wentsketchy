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
	batches batches,
	fifoPath string,
) (batches, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("calendar: could not generate update event")
	}

	calendarItem := sketchybar.ItemOptions{
		Padding: sketchybar.PaddingOptions{
			Left:  settings.SketchybarSettings.ItemSpacing,
			Right: settings.SketchybarSettings.ItemSpacing,
		},
		Icon: sketchybar.ItemIconOptions{
			Value: settings.IconClock,
			Padding: sketchybar.PaddingOptions{
				Left:  settings.SketchybarSettings.IconPadding,
				Right: settings.SketchybarSettings.IconPadding / 2,
			},
		},
		Background: sketchybar.BackgroundOptions{
			CornerRadius: settings.SketchybarSettings.ItemRadius,
		},
		Label: sketchybar.ItemLabelOptions{
			Padding: sketchybar.PaddingOptions{
				Left:  0,
				Right: settings.SketchybarSettings.IconPadding,
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
	batches batches,
	args *args.In,
) (batches, error) {
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
