package aerospace

import (
	"strings"
)

const (
	outputFormatAppBundleID  = "%{app-bundle-id}"
	outputFormatAppName      = "%{app-name}"
	outputFormatAppPid       = "%{app-pid}"
	outputFormatTab          = "%{tab}"
	outputFormatWindowID     = "%{window-id}"
	outputFormatWindowTitle  = "%{window-title}"
	outputFormatWorkspace    = "%{workspace}"
	outputFormatMonitorID    = "%{monitor-id}"
	outputFormatMonitorName  = "%{monitor-name}"
	outputFormatRightPadding = "%{right-padding}"
	outputFormatNewline      = "%{newline}"
)

const outputFormatDefaultApp = "%{app-pid}%{right-padding} | %{app-bundle-id}%{right-padding} | %{app-name}"
const outputFormatDefaultWindow = "%{window-id}%{right-padding} | %{app-name}%{right-padding} | %{window-title}"
const outputFormatDefaultWorkspace = "%{monitor-id}%{right-padding} | %{monitor-name}"
const outputFormatDefaultMonitor = "%{workspace}"

const outputFormatSeparator = "¬"

func windowOutputFormat() string {
	return strings.Join(
		[]string{
			outputFormatWindowID,
			outputFormatSeparator,
			outputFormatAppName,
			outputFormatSeparator,
		}, "",
	)
}

func fullWindowOutputFormat() string {
	return strings.Join(
		[]string{
			outputFormatWindowID,
			outputFormatSeparator,
			outputFormatAppName,
			outputFormatSeparator,
			outputFormatWorkspace,
			outputFormatSeparator,
			outputFormatMonitorID,
		}, "",
	)
}

func fullWorkspaceOutputFormat() string {
	return strings.Join(
		[]string{
			outputFormatWorkspace,
			outputFormatSeparator,
			outputFormatMonitorID,
		}, "",
	)
}

func workspaceOutputFormat() string {
	return outputFormatWorkspace
}

func monitorOutputFormat() string {
	return outputFormatMonitorID
}
