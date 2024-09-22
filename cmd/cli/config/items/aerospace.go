package items

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
	"github.com/lucax88x/wentsketchy/internal/utils"
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
	batches batches,
	fifoPath string,
) (batches, error) {
	batches, err := i.applyTree(batches)

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
	batches batches,
	args *args.In,
) (batches, error) {
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
		Padding: sketchybar.PaddingOptions{
			Left:  settings.SketchybarSettings.ItemSpacing,
			Right: 0,
		},
		Background: sketchybar.BackgroundOptions{
			CornerRadius: 0,
			// BorderOptions: sketchybar.BorderOptions{
			// 	Color: colors-Red,
			// 	Width: 10,
			// },
		},
		Icon: sketchybar.ItemIconOptions{
			Value: icon,
			Color: sketchybar.ColorOptions{
				Color: colors.White,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.SketchybarSettings.IconPadding,
				Right: settings.SketchybarSettings.IconPadding,
			},
		},
		ClickScript: fmt.Sprintf(`aerospace workspace "%s"`, workspaceID),
	}, nil
}

func (i *AerospaceItem) windowToSketchybar(
	monitorID int,
	workspaceID string,
	window *aerospace.Window,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := settings.AppIcons[window.App]

	if !hasIcon {
		return nil, fmt.Errorf("could not find icon for app %s", window.App)
	}

	itemOptions := &sketchybar.ItemOptions{
		Display: strconv.Itoa(monitorID),
		Background: sketchybar.BackgroundOptions{
			CornerRadius: 0,
		},
		Icon: sketchybar.ItemIconOptions{
			Color: sketchybar.ColorOptions{
				Color:          colors.White,
				HighlightColor: colors.White,
			},
			Font: sketchybar.FontOptions{
				Font: settings.FontAppIcon,
				Kind: "Regular",
				Size: "14.0",
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.SketchybarSettings.IconPadding,
				Right: settings.SketchybarSettings.IconPadding,
			},
			Value: icon,
		},
		ClickScript: fmt.Sprintf(`aerospace workspace "%s"`, workspaceID),
	}

	if utils.Equals(window.App, i.aerospace.GetFocusedApp()) {
		itemOptions.Background.Color = sketchybar.ColorOptions{
			Color: settings.SketchybarSettings.AerospaceItemFocusedBackgroundColor,
		}
	}

	return itemOptions, nil
}

func sketchybarSpaceID(spaceID string) string {
	return fmt.Sprintf("%s.%s", spaceItemPrefix, spaceID)
}

func sketchybarWindowID(window *aerospace.Window) string {
	return fmt.Sprintf("%s.%d", windowItemPrefix, window.ID)
}

func (i *AerospaceItem) addWindowToSketchybar(
	batches batches,
	monitorID int,
	workspaceID string,
	window *aerospace.Window,
) (batches, error) {
	windowItem, err := i.windowToSketchybar(
		monitorID,
		workspaceID,
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

func checker(batches batches, fifoPath string) (batches, error) {
	updateEvent, err := args.BuildEvent(fifoPath)

	if err != nil {
		return batches, errors.New("aerospace: could not generate update event")
	}

	checkerItem := sketchybar.ItemOptions{
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
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
	batches batches,
) (batches, error) {
	indexedWindows, err := i.aerospace.WindowsOfFocusedMonitor(ctx)

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
	focusedSpaceID := i.aerospace.GetFocusedWorkspaceID(ctx)
	focusedSketchybarSpaceID := sketchybarSpaceID(focusedSpaceID)

	for _, window := range indexedWindows {
		_, found := alreadyThereWindows[window.ID]
		if !found {
			batches, err = i.addWindowToSketchybar(batches, focusedMonitorID, focusedSpaceID, window)

			windowID := sketchybarWindowID(window)

			batches = batch(batches, s("--move", windowID, "after", focusedSketchybarSpaceID))

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

func (i AerospaceItem) handleSpaceWindowsChange2(ctx context.Context, batches batches) (batches, error) {
	i.aerospace.RefreshTree()

	batches = batch(batches, s("--remove", fmt.Sprintf("/%s/", spaceItemPrefix)))
	batches = batch(batches, s("--remove", fmt.Sprintf("/%s/", windowItemPrefix)))

	batches, err := i.applyTree(batches)

	if err != nil {
		return batches, err
	}

	return batches, nil
}

func (i *AerospaceItem) applyTree(batches batches) (batches, error) {
	tree := i.aerospace.GetTree()

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

				batches, err = i.addWindowToSketchybar(
					batches,
					monitor.Monitor,
					workspace.Workspace,
					window,
				)

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
