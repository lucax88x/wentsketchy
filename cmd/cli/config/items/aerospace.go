package items

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
)

type AerospaceItem struct {
	aerospace  aerospace.Aerospace
	sketchybar sketchybar.API
}

func NewAerospaceItem(
	aerospace aerospace.Aerospace,
	sketchybar sketchybar.API,
) AerospaceItem {
	return AerospaceItem{
		aerospace,
		sketchybar,
	}
}

const aerospaceCheckerItemName = "aerospace.checker"
const spaceItemPrefix = "workspace"
const windowItemPrefix = "window"

func (i AerospaceItem) Init(
	batches [][]string,
	fifoPath string,
) ([][]string, error) {
	batches, err := applyTree(batches, i.aerospace.GetTree())

	if err != nil {
		return batches, err
	}

	batches, err = checker(batches, fifoPath)

	if err != nil {
		return batches, err
	}

	return batches, nil
}

func (i AerospaceItem) Update(
	ctx context.Context,
	batches [][]string,
	args *args.In,
) ([][]string, error) {
	if !isAerospace(args.Name) {
		return batches, nil
	}

	var err error
	if args.Event == events.SpaceWindowsChange {
		batches, err = i.handleSpaceWindowsChange2(ctx, batches)

		if err != nil {
			return batches, err
		}
	}

	if args.Event == events.DisplayChange {
		i.aerospace.SetFocusedWorkspaceID(args.Info)
	}

	return batches, err
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
		// Space:   workspaceID,
		Display: strconv.Itoa(monitorID),
		Background: sketchybar.BackgroundOptions{
			Drawing:      true,
			ColorOptions: sketchybar.ColorOptions{Color: settings.ColorBlue},
			CornerRadius: 10,
		},
		// Border: sketchybar.BorderOptions{
		// 	Color:  settings.ColorBlue,
		// 	Width:  10,
		// 	Height: 10,
		// },
		Icon: sketchybar.ItemIconOptions{
			Value: icon,
			Font: sketchybar.FontOptions{
				Font: settings.FontIcon,
				Kind: "Regular",
				Size: "12.0",
			},
			PaddingOptions: sketchybar.PaddingOptions{
				Right: 4,
				Left:  4,
			},
		},
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
				Right: 4,
				Left:  4,
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

func sketchybarSpaceID(spaceID string) string {
	return fmt.Sprintf("%s.%s", spaceItemPrefix, spaceID)
}

func sketchybarWindowID(window *aerospace.Window) string {
	return fmt.Sprintf("%s.%d", windowItemPrefix, window.ID)
}

func addWindowToSketchybar(
	batches [][]string,
	monitorID int,
	window *aerospace.Window,
) ([][]string, error) {
	windowItem, err := windowToSketchybar(
		monitorID,
		window,
	)

	if err != nil {
		return batches, err
	}

	windowID := sketchybarWindowID(window)
	batches = batch(batches, s("--add", "item", windowID, "left"))
	batches = batch(batches, m(s("--set", windowID), windowItem.ToArgs()))
	return batches, err
}

func checker(batches [][]string, fifoPath string) ([][]string, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("aerospace: could not generate update event")
	}

	checkerItem := sketchybar.ItemOptions{
		Background: sketchybar.BackgroundOptions{
			Drawing: false,
		},
		Updates: "on",
		Script:  updateEvent,
	}

	batches = batch(batches, s("--add", "item", aerospaceCheckerItemName, "left"))
	batches = batch(batches, m(s("--set", aerospaceCheckerItemName), checkerItem.ToArgs()))
	batches = batch(batches, s("--subscribe", aerospaceCheckerItemName,
		events.DisplayChange,
		events.SpaceWindowsChange,
		events.SystemWoke,
	))

	return batches, nil
}

func (i AerospaceItem) handleSpaceWindowsChange(
	ctx context.Context,
	batches [][]string,
) ([][]string, error) {
	indexedWindows, err := i.aerospace.WindowsOfFocusedMonitor()

	if err != nil {
		return batches, fmt.Errorf("aerospace: cannot query aerospace windows. %w", err)
	}

	actualBar, err := i.sketchybar.QueryBar(ctx)

	if err != nil {
		return batches, fmt.Errorf("aerospace: cannot query sketchybar bar. %w", err)
	}

	alreadyThereWindows := make(map[int]bool, 0)
	for _, barItemID := range actualBar.Items {
		if strings.HasPrefix(barItemID, windowItemPrefix) {
			withoutPrefix, _ := strings.CutPrefix(barItemID, fmt.Sprintf("%s.", windowItemPrefix))

			converted, err := strconv.Atoi(withoutPrefix)

			if err != nil {
				return batches, fmt.Errorf(
					"aerospace: cannot convert id '%s' to int. %w",
					withoutPrefix,
					err,
				)
			}

			_, foundWindow := indexedWindows[converted]

			if !foundWindow {
				batches = batch(batches, s("--remove", barItemID))
			} else {
				alreadyThereWindows[converted] = true
			}
		}
	}

	focusedMonitorID := i.aerospace.GetFocusedMonitorID(ctx)
	focusedSpaceID := sketchybarSpaceID(i.aerospace.GetFocusedWorkspaceID(ctx))

	for _, window := range indexedWindows {
		_, found := alreadyThereWindows[window.ID]
		if !found {
			batches, err = addWindowToSketchybar(batches, focusedMonitorID, window)

			windowID := sketchybarWindowID(window)

			batches = batch(batches, s("--move", windowID, "after", focusedSpaceID))

			if err != nil {
				return batches, fmt.Errorf(
					"aerospace: cannot add window '%d'. %w",
					window.ID,
					err,
				)
			}
		}
	}

	return batches, nil
}

func (i AerospaceItem) handleSpaceWindowsChange2(
	ctx context.Context,
	batches [][]string,
) ([][]string, error) {
	i.aerospace.RefreshTree()

	batches = batch(batches, s("--remove", fmt.Sprintf("/%s/", spaceItemPrefix)))
	batches = batch(batches, s("--remove", fmt.Sprintf("/%s/", windowItemPrefix)))

	batches, err := applyTree(batches, i.aerospace.GetTree())

	if err != nil {
		return batches, err
	}

	return batches, nil
}

func applyTree(batches [][]string, tree *aerospace.Tree) ([][]string, error) {
	var aggregatedErr error
	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			spaceID := sketchybarSpaceID(workspace.Workspace)
			workspaceSpace, err := workspaceToSketchybar(monitor.Monitor, workspace.Workspace)

			if err != nil {
				aggregatedErr = errors.Join(aggregatedErr, err)
				continue
			}

			batches = batch(batches, s("--add", "item", spaceID, "left"))
			batches = batch(batches, m(s("--set", spaceID), workspaceSpace.ToArgs()))

			for _, windowID := range workspace.Windows {
				window := tree.IndexedWindows[windowID]

				batches, err = addWindowToSketchybar(batches, monitor.Monitor, window)

				if err != nil {
					aggregatedErr = errors.Join(aggregatedErr, err)
					continue
				}
			}
		}
	}

	return batches, aggregatedErr
}

func isAerospace(name string) bool {
	return name == aerospaceCheckerItemName
}
