package args

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/fifo"
)

// https://felixkratz.github.io/SketchyBar/config/events
type In struct {
	// the item name
	Name string `json:"name"`
	// the event
	Event    string `json:"event"`
	Info     string `json:"info"`
	Button   string `json:"button"`
	Modifier string `json:"modifier"`
}

// $INFO is a json, and its not easy to embed a json inside a json
type Out struct {
	Name     string `json:"name"`
	Event    string `json:"event"`
	Button   string `json:"button"`
	Modifier string `json:"modifier"`
}

func FromEvent(msg string) (*In, error) {
	// Find the positions of the fixed parts
	argsStart := strings.Index(msg, "args: ") + len("args: ")
	infoStart := strings.Index(msg, "info: ") + len("info: ")

	// Extract the JSON substrings
	argsJSON := msg[argsStart : infoStart-len(" info: ")]
	infoJSON := msg[infoStart:]

	var args *In
	err := json.Unmarshal([]byte(argsJSON), &args)

	if err != nil {
		return nil, fmt.Errorf("args: could not deserialize data. %w", err)
	}

	args.Info = infoJSON

	return args, nil
}

func BuildEvent() (string, error) {
	data := &Out{
		Name:     "$NAME",
		Event:    "$SENDER",
		Button:   "$BUTTON",
		Modifier: "$MODIFIER",
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		return "", fmt.Errorf("args: could not serialize data. %w", err)
	}

	serialized := strings.ReplaceAll(string(bytes), `"`, `\"`)

	// TODO: ensure file exists, also in aerospace.toml
	return fmt.Sprintf(
		// `[[ -f %s ]] && echo "update args: %s info: $INFO %c" >> %s`,
		`echo "update args: %s info: $INFO %c" >> %s`,
		// settings.FifoPath,
		serialized,
		fifo.Separator,
		settings.FifoPath,
	), nil
}
