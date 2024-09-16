package settings

const (
	IconApple           = ""
	IconClock           = ""
	IconChat            = "󱅱"
	IconTerminal        = ""
	IconCode            = ""
	IconChrome          = ""
	IconFinder          = ""
	IconEmail           = "󰇰"
	IconTools           = ""
	IconDocuments       = "󰧮"
	IconBattery100      = "􀛨"
	IconBattery75       = "􀺸"
	IconBattery50       = "􀺶"
	IconBattery25       = "􀛩"
	IconBattery0        = "􀛪"
	IconBatteryCharging = "􀢋"
)

//nolint:gochecknoglobals // ok
var WorkspaceIcons = map[string]string{
	"1": IconChat,
	"2": IconCode,
	"3": IconTerminal,
	"4": IconChrome,
	"5": IconFinder,
	"6": IconEmail,
	"7": IconTools,
	"8": IconDocuments,
}
