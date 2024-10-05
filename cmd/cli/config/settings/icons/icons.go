package icons

const (
	Apple           = ""
	Clock           = ""
	Chat            = "󱅱"
	Terminal        = ""
	Code            = ""
	Chrome          = ""
	Finder          = ""
	Email           = "󰇰"
	Tools           = ""
	CPU             = "􀫥"
	ThermoMedium    = "􀇬"
	Documents       = "󰧮"
	Battery100      = "􀛨"
	Battery75       = "􀺸"
	Battery50       = "􀺶"
	Battery25       = "􀛩"
	Battery0        = "􀛪"
	BatteryCharging = "􀢋"
)

//nolint:gochecknoglobals // ok
var Workspace = map[string]string{
	"1": Chat,
	"2": Code,
	"3": Terminal,
	"4": Chrome,
	"5": Finder,
	"6": Email,
	"7": Tools,
	"8": Documents,
}
