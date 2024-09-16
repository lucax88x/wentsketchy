package args

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// https://felixkratz.github.io/SketchyBar/config/events
type Args struct {
	// the item name
	Name string `json:"name"`
	// the event
	Event    string `json:"event"`
	Info     string `json:"info"`
	Button   string `json:"button"`
	Modifier string `json:"modifier"`
}

func FromMsg(msg string) (*Args, error) {
	msg, _ = strings.CutPrefix(msg, "update")

	var args *Args
	err := json.Unmarshal([]byte(msg), &args)

	if err != nil {
		return nil, fmt.Errorf("args: could not deserialize data. %w", err)
	}

	return args, nil
}

func BuildEvent(path string) (string, error) {
	if path == "" {
		return "", errors.New("args: path is empty")
	}

	data := &Args{
		Name:     "$NAME",
		Event:    "$SENDER",
		Info:     "$INFO",
		Button:   "$BUTTON",
		Modifier: "$MODIFIER",
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		return "", fmt.Errorf("args: could not serialize data. %w", err)
	}

	serialized := strings.ReplaceAll(string(bytes), `"`, `\"`)

	return fmt.Sprintf(`echo "update %s" > %s`, serialized, path), nil
}
