package config

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/wentsketchy"
)

func initAerospace(ctx context.Context, di *wentsketchy.Wentsketchy) ([]string, error) {
	tree, err := di.Aerospace.Build(ctx)

	if err != nil {
		return make([]string, 0), err
	}

	var batches = make([][]string, 0)
	var aggregatedErr error
	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			spaceID := fmt.Sprintf("space.%s", workspace.Workspace.Id)
			workspaceSpace, err := workspaceToSketchybar(monitor.Monitor, workspace.Workspace)

			if err != nil {
				aggregatedErr = errors.Join(aggregatedErr, err)
				continue
			}

			batches = batch(batches, s("--add", "item", spaceID, "left"))
			batches = batch(batches, m(s("--set", spaceID), workspaceSpace.ToArgs()))

			for _, window := range workspace.Windows {
				windowItem, err := windowToSketchybar(
					monitor.Monitor,
					window,
				)

				if err != nil {
					aggregatedErr = errors.Join(aggregatedErr, err)
					continue
				}

				windowID := fmt.Sprintf("window.%s", window.Id)
				batches = batch(batches, s("--add", "item", windowID, "left"))
				batches = batch(batches, m(s("--set", windowID), windowItem.ToArgs()))
			}
		}
	}

	return flatten(batches...), nil
}

func workspaceToSketchybar(
	monitor *aerospace.Monitor,
	workspace *aerospace.Workspace,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := wentsketchy.WorkspaceIcons[workspace.Id]

	if !hasIcon {
		return nil, fmt.Errorf("could not find icon for workspace %s", workspace.Id)
	}

	return &sketchybar.ItemOptions{
		// Space:   workspace.Id,
		Display: strconv.Itoa(monitor.Id),
		//   padding_left=1
		//   padding_right=1
		Icon: sketchybar.ItemIconOptions{
			//   icon.highlight_color="$RED"
			Value: icon,
			Font: sketchybar.FontOptions{
				Font: wentsketchy.FontIcon,
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
				Color: wentsketchy.ColorBackground2,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: wentsketchy.ColorBackground1,
			},
		},
		//   script="$PLUGIN_DIR/space.sh"
	}, nil
}

func windowToSketchybar(
	monitor *aerospace.Monitor,
	window *aerospace.Window,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := wentsketchy.AppIcons[window.App]

	if !hasIcon {
		return nil, fmt.Errorf("could not find icon for app %s", window.App)
	}

	return &sketchybar.ItemOptions{
		Display: strconv.Itoa(monitor.Id),
		//   padding_left=1
		//   padding_right=1
		Icon: sketchybar.ItemIconOptions{
			ColorOptions: sketchybar.ColorOptions{
				Color:          wentsketchy.ColorWhite,
				HighlightColor: wentsketchy.ColorRed,
			},
			//   icon.highlight_color="$RED"
			Value: icon,
			Font: sketchybar.FontOptions{
				Font: wentsketchy.FontAppIcon,
				Kind: "Regular",
				Size: "14.0",
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
				Color: wentsketchy.ColorBackground2,
			},
			ColorOptions: sketchybar.ColorOptions{
				Color: wentsketchy.ColorBackground1,
			},
		},
		//   script="$PLUGIN_DIR/space.sh"
	}, nil
}
