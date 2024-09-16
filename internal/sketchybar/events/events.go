package events

import (
	"strings"
)

type Event string

const (
	FrontAppSwitched    Event = "front_app_switched"
	SpaceWindowsChange  Event = "space_windows_change"
	SpaceChange         Event = "space_change"
	DisplayChange       Event = "display_change"
	VolumeChange        Event = "volume_change"
	BrightnessChange    Event = "brightness_change"
	PowerSourceChanged  Event = "power_source_change"
	WifiChange          Event = "wifi_change"
	MediaChange         Event = "media_change"
	SystemWillSleep     Event = "system_will_sleep"
	SystemWoke          Event = "system_woke"
	MouseEntered        Event = "mouse.entered"
	MouseExited         Event = "mouse.exited"
	MouseEnteredGlobal  Event = "mouse.entered.global"
	MouseExitedGlobal   Event = "mouse.exited.global"
	MouseClicked        Event = "mouse.clicked"
	MouseScrolled       Event = "mouse.scrolled"
	MouseScrolledGlobal Event = "mouse.scrolled.global"
)

func ToString(events ...Event) string {
	casted := make([]string, len(events))
	for i, event := range events {
		casted[i] = string(event)
	}
	return strings.Join(casted, " ")
}
