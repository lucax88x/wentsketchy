package settings

type Settings struct {
	LabelColor    string
	LabelFont     string
	LabelFontKind string
	LabelFontSize string
	IconColor     string
	IconFont      string
	IconFontKind  string
	IconFontSize  string
	IconStripFont string
}

var SketchybarSettings = Settings{
	LabelColor:    ColorWhite,
	LabelFont:     FontLabel,
	LabelFontKind: "Semibold",
	LabelFontSize: "14.0",
	IconColor:     ColorWhite,
	IconFont:      FontIcon,
	IconFontKind:  "Regular",
	IconFontSize:  "14.0",
	IconStripFont: FontAppIcon,
}
