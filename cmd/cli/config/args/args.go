package args

import (
	"fmt"
	"strings"
)

// https://felixkratz.github.io/SketchyBar/config/events
type Args struct {
	// the item name
	Name string
	// the event
	Event    string
	Info     string
	Button   string
	Modifier string
}

func FromMsg(msg string) *Args {
	msg = strings.Replace(msg, "update ", "", 1)

	args := strings.Split(msg, "|")

	name := ""
	event := ""
	info := ""
	button := ""
	modifier := ""

	if len(args) > 0 {
		name = args[0]
	}

	if len(args) > 1 {
		event = args[1]
	}

	if len(args) > 2 {
		info = args[2]
	}

	if len(args) > 3 {
		button = args[3]
	}

	if len(args) > 4 {
		modifier = args[4]
	}

	return &Args{
		Name:     name,
		Event:    event,
		Info:     info,
		Button:   button,
		Modifier: modifier,
	}
}

func BuildEvent(path string) string {
	return fmt.Sprintf(`echo "update $NAME|$SENDER|$INFO|$BUTTON|$MODIFIER" > %s`, path)
}
