package items

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
)

type AerospaceItem struct {
	AerospaceData *aerospace.Data
}

func NewAerospaceItem(data *aerospace.Data) AerospaceItem {
	return AerospaceItem{data}
}

func (i AerospaceItem) Init(
	batches [][]string,
) ([][]string, error) {
	tree := i.AerospaceData.Tree

	var aggregatedErr error
	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			spaceID := fmt.Sprintf("space.%s", workspace.Workspace)
			workspaceSpace, err := workspaceToSketchybar(monitor.Monitor, workspace.Workspace)

			if err != nil {
				aggregatedErr = errors.Join(aggregatedErr, err)
				continue
			}

			batches = batch(batches, s("--add", "item", spaceID, "left"))
			batches = batch(batches, m(s("--set", spaceID), workspaceSpace.ToArgs()))

			for _, windowID := range workspace.Windows {
				window := tree.IndexedWindows[windowID]

				windowItem, err := windowToSketchybar(
					monitor.Monitor,
					window,
				)

				if err != nil {
					aggregatedErr = errors.Join(aggregatedErr, err)
					continue
				}

				windowID := fmt.Sprintf("window.%d", windowID)
				batches = batch(batches, s("--add", "item", windowID, "left"))
				batches = batch(batches, m(s("--set", windowID), windowItem.ToArgs()))
			}
		}
	}

	return batches, nil
}

func workspaceToSketchybar(
	monitorID int,
	workspaceID string,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := settings.WorkspaceIcons[workspaceID]

	if !hasIcon {
		return nil, fmt.Errorf("could not find icon for workspace %s", workspaceID)
	}

	return &sketchybar.ItemOptions{
		// Space:   workspace.Id,
		Display: strconv.Itoa(monitorID),
		//   padding_left=1
		//   padding_right=1
		Icon: sketchybar.ItemIconOptions{
			//   icon.highlight_color="$RED"
			Value: icon,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
				Kind: "Regular",
				Size: "12.0",
			},
			//   icon.highlight_color="$RED"
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
		},
		Background: sketchybar.BackgroundOptions{
			//   background.drawing="on"
			BorderOptions: sketchybar.BorderOptions{
				Color: settings.ColorBackground2,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: settings.ColorBackground1,
			},
		},
		//   script="$PLUGIN_DIR/space.sh"
	}, nil
}

func windowToSketchybar(
	monitorID int,
	window *aerospace.Window,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := settings.AppIcons[window.App]

	if !hasIcon {
		return nil, fmt.Errorf("could not find icon for app %s", window.App)
	}

	return &sketchybar.ItemOptions{
		Display: strconv.Itoa(monitorID),
		//   padding_left=1
		//   padding_right=1
		Icon: sketchybar.ItemIconOptions{
			ColorOptions: sketchybar.ColorOptions{
				Color:          settings.ColorWhite,
				HighlightColor: settings.ColorRed,
			},
			Font: sketchybar.FontOptions{
				Font: settings.FontAppIcon,
				Kind: "Regular",
				Size: "14.0",
			},
			//   icon.highlight_color="$RED"
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 5,
				Left:  5,
			},
			Value:     icon,
			Highlight: false,
		},
		Background: sketchybar.BackgroundOptions{
			//   background.drawing="on"
			BorderOptions: sketchybar.BorderOptions{
				Color: settings.ColorBackground2,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: settings.ColorBackground1,
			},
		},
		//   script="$PLUGIN_DIR/space.sh"
	}, nil
}
