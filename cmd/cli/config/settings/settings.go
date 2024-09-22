package settings

import "github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"

type Settings struct {
	BarBackgroundColor                  string
	BarHeight                           int
	BarMargin                           int
	ItemHeight                          int
	ItemSpacing                         int
	ItemRadius                          int
	ItemBackgroundColor                 string
	AerospaceItemFocusedBackgroundColor string
	IconPadding                         int
	LabelColor                          string
	LabelFont                           string
	LabelFontKind                       string
	LabelFontSize                       string
	IconColor                           string
	IconFont                            string
	IconFontKind                        string
	IconFontSize                        string
	IconStripFont                       string
}

var SketchybarSettings = Settings{
	BarBackgroundColor:                  colors.Transparent,
	BarHeight:                           35,
	BarMargin:                           4,
	ItemHeight:                          25,
	ItemSpacing:                         12,
	ItemRadius:                          45,
	IconPadding:                         12,
	ItemBackgroundColor:                 colors.Black1,
	AerospaceItemFocusedBackgroundColor: colors.Black3,
	LabelColor:                          colors.White,
	LabelFont:                           FontLabel,
	LabelFontKind:                       "Semibold",
	LabelFontSize:                       "14.0",
	IconColor:                           colors.White,
	IconFont:                            FontIcon,
	IconFontKind:                        "Bold",
	IconFontSize:                        "16.0",
	IconStripFont:                       FontAppIcon,
}
