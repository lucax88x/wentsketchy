package events

// https://felixkratz.github.io/SketchyBar/config/events
const (
	Forced              string = "forced"
	Routine             string = "routine"
	FrontAppSwitched    string = "front_app_switched"
	SpaceWindowsChange  string = "space_windows_change"
	SpaceChange         string = "space_change"
	DisplayChange       string = "display_change"
	VolumeChange        string = "volume_change"
	BrightnessChange    string = "brightness_change"
	PowerSourceChanged  string = "power_source_change"
	WifiChange          string = "wifi_change"
	MediaChange         string = "media_change"
	SystemWillSleep     string = "system_will_sleep"
	SystemWoke          string = "system_woke"
	MouseEntered        string = "mouse.entered"
	MouseExited         string = "mouse.exited"
	MouseEnteredGlobal  string = "mouse.entered.global"
	MouseExitedGlobal   string = "mouse.exited.global"
	MouseClicked        string = "mouse.clicked"
	MouseScrolled       string = "mouse.scrolled"
	MouseScrolledGlobal string = "mouse.scrolled.global"
)
