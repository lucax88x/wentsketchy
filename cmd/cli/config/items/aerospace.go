package items

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/colors"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/settings/icons"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	aerospace_events "github.com/lucax88x/wentsketchy/internal/aerospace/events"
	"github.com/lucax88x/wentsketchy/internal/sketchybar"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/events"
	"github.com/lucax88x/wentsketchy/internal/utils"
)

type AerospaceItem struct {
	logger     *slog.Logger
	aerospace  aerospace.Aerospace
	sketchybar sketchybar.API

	position  sketchybar.Position
	windowIDs map[aerospace.WindowID]aerospace.WorkspaceID
}

func NewAerospaceItem(
	logger *slog.Logger,
	aerospace aerospace.Aerospace,
	sketchybarAPI sketchybar.API,
) AerospaceItem {
	return AerospaceItem{
		logger,
		aerospace,
		sketchybarAPI,
		sketchybar.PositionLeft,
		make(map[int]string, 0),
	}
}

const aerospaceCheckerItemName = "aerospace.checker"
const workspaceItemPrefix = "aerospace.workspace"
const windowItemPrefix = "aerospace.window"
const bracketItemPrefix = "aerospace.bracket"
const spacerItemPrefix = "aerospace.spacer"

const AerospaceName = aerospaceCheckerItemName

func (item AerospaceItem) Init(
	ctx context.Context,
	position sketchybar.Position,
	batches Batches,
) (Batches, error) {
	item.position = position

	batches, err := item.applyTree(ctx, batches, position)

	if err != nil {
		return batches, err
	}

	batches, err = checker(batches, position)

	if err != nil {
		return batches, err
	}

	return batches, nil
}

func (item AerospaceItem) Update(
	ctx context.Context,
	batches Batches,
	position sketchybar.Position,
	args *args.In,
) (Batches, error) {
	item.position = position

	if !isAerospace(args.Name) {
		return batches, nil
	}

	var err error
	if args.Event == aerospace_events.WorkspaceChange {
		var data aerospace_events.WorkspaceChangeEventInfo
		err := json.Unmarshal([]byte(args.Info), &data)

		if err != nil {
			return batches, fmt.Errorf("aerospace: could not deserialize json for workspace-change. %v", args.Info)
		}

		batches = item.handleWorkspaceChange(ctx, batches, data.Prev, data.Focused)
	}

	if args.Event == events.SpaceWindowsChange {
		batches, err = item.CheckTree(ctx, batches)

		if err != nil {
			return batches, err
		}
	}

	if args.Event == events.DisplayChange {
		batches = item.handleDisplayChange(batches)
	}

	if args.Event == events.FrontAppSwitched {
		batches = item.handleFrontAppSwitched(ctx, batches, args.Info)
	}

	return batches, err
}

func (item AerospaceItem) CheckTree(
	ctx context.Context,
	batches Batches,
) (Batches, error) {
	item.aerospace.SingleFlightRefreshTree()
	tree := item.aerospace.GetTree()

	focusedWorkspaceID := item.aerospace.GetFocusedWorkspaceID(ctx)

	var aggregatedErr error
	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			isFocusedWorkspace := workspace.Workspace == focusedWorkspaceID

			for _, windowID := range workspace.Windows {
				window := tree.IndexedWindows[windowID]

				currentWorkspaceID, found := item.windowIDs[windowID]

				if found {
					if currentWorkspaceID != workspace.Workspace {
						item.logger.InfoContext(
							ctx,
							"aerospace: moving item to workspace",
							slog.Int("window.id", window.ID),
							slog.String("window.app", window.App),
							slog.String("workspace.id", workspace.Workspace),
						)

						batches = batch(batches, s(
							"--move",
							getSketchybarWindowID(windowID),
							"after",
							getSketchybarWorkspaceID(workspace.Workspace),
						))
					}
				} else {
					var sketchybarWindowID string
					batches, sketchybarWindowID = item.addWindowToSketchybar(
						batches,
						item.position,
						isFocusedWorkspace,
						monitor.Monitor,
						workspace.Workspace,
						window.ID,
						window.App,
					)

					item.logger.InfoContext(
						ctx,
						"aerospace: added item to workspace",
						slog.Int("window.id", window.ID),
						slog.String("window.app", window.App),
						slog.String("workspace.id", workspace.Workspace),
						slog.String("sketchybar", sketchybarWindowID),
					)

					batches = batch(batches, s(
						"--move",
						sketchybarWindowID,
						"after",
						getSketchybarWorkspaceID(workspace.Workspace),
					))
				}
			}
		}
	}

	for windowID := range item.windowIDs {
		_, found := tree.IndexedWindows[windowID]

		if !found {
			batches = item.removeWindow(batches, windowID)
		}
	}

	return batches, aggregatedErr
}

func (item AerospaceItem) handleFrontAppSwitched(
	ctx context.Context,
	batches Batches,
	windowApp string,
) Batches {
	item.aerospace.SetFocusedApp(windowApp)

	focusedWorkspaceID := item.aerospace.GetFocusedWorkspaceID(ctx)

	batches = item.highlightWindows(batches, focusedWorkspaceID, windowApp)

	return batches
}

func (item AerospaceItem) workspaceToSketchybar(
	isFocusedWorkspace bool,
	monitorsCount int,
	monitorID int,
	workspaceID string,
) (*sketchybar.ItemOptions, error) {
	icon, hasIcon := icons.Workspace[workspaceID]

	if !hasIcon {
		item.logger.Info(
			"could not find icon for app",
			slog.String("app", workspaceID),
		)
		return nil, fmt.Errorf("could not find icon for workspace %s", workspaceID)
	}

	colors := item.getWorkspaceColors(isFocusedWorkspace)

	return &sketchybar.ItemOptions{
		Display: item.getSketchybarDisplayIndex(monitorsCount, monitorID),
		Padding: sketchybar.PaddingOptions{
			Left:  pointer(0),
			Right: pointer(0),
		},
		Background: sketchybar.BackgroundOptions{
			Drawing: "on",
			Color: sketchybar.ColorOptions{
				Color: colors.backgroundColor,
			},
		},
		Icon: sketchybar.ItemIconOptions{
			Value: icon,
			Color: sketchybar.ColorOptions{
				Color: colors.color,
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.Sketchybar.Aerospace.Padding,
				Right: settings.Sketchybar.Aerospace.Padding,
			},
		},
		ClickScript: fmt.Sprintf(`aerospace workspace "%s"`, workspaceID),
	}, nil
}

func (item *AerospaceItem) windowToSketchybar(
	isFocusedWorkspace bool,
	monitorID aerospace.MonitorID,
	workspaceID aerospace.WorkspaceID,
	windowApp string,
) *sketchybar.ItemOptions {
	icon, hasIcon := icons.App[windowApp]

	if !hasIcon {
		item.logger.Info(
			"could not find icon for app",
			slog.String("app", windowApp),
		)
		icon = ""
	}

	windowVisibility := item.getWindowVisibility(isFocusedWorkspace)
	itemOptions := &sketchybar.ItemOptions{
		Display: strconv.Itoa(monitorID),
		Width:   windowVisibility.width,
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
		Icon: sketchybar.ItemIconOptions{
			Drawing: windowVisibility.show,
			Color: sketchybar.ColorOptions{
				Color: windowVisibility.color,
			},
			Font: sketchybar.FontOptions{
				Font: settings.FontAppIcon,
				Kind: "Regular",
				Size: "14.0",
			},
			Padding: sketchybar.PaddingOptions{
				Left:  settings.Sketchybar.Aerospace.Padding,
				Right: settings.Sketchybar.Aerospace.Padding,
			},
			Value: icon,
		},
		ClickScript: fmt.Sprintf(`aerospace workspace "%s"`, workspaceID),
	}

	if utils.Equals(windowApp, item.aerospace.GetFocusedApp()) {
		itemOptions.Icon.Color = sketchybar.ColorOptions{
			Color: windowVisibility.focusedColor,
		}
	}

	return itemOptions
}

func getSketchybarWorkspaceID(spaceID aerospace.WorkspaceID) string {
	return fmt.Sprintf("%s.%s", workspaceItemPrefix, spaceID)
}

func getSketchybarWindowID(windowID aerospace.WindowID) string {
	return fmt.Sprintf("%s.%d", windowItemPrefix, windowID)
}

func getSketchybarBracketID(spaceID aerospace.WorkspaceID) string {
	return fmt.Sprintf("%s.%s", bracketItemPrefix, spaceID)
}

func getSketchybarSpacerID(spaceID aerospace.WorkspaceID) string {
	return fmt.Sprintf("%s.%s", spacerItemPrefix, spaceID)
}

func (item *AerospaceItem) removeWindow(batches Batches, windowID aerospace.WindowID) Batches {
	delete(item.windowIDs, windowID)

	batches = batch(batches, s("--remove", getSketchybarWindowID(windowID)))

	return batches
}

func (item *AerospaceItem) addWindowToSketchybar(
	batches Batches,
	position sketchybar.Position,
	isFocusedWorkspace bool,
	monitorID aerospace.MonitorID,
	workspaceID aerospace.WorkspaceID,
	windowID aerospace.WindowID,
	windowApp string,
) (Batches, string) {
	windowItem := item.windowToSketchybar(
		isFocusedWorkspace,
		monitorID,
		workspaceID,
		windowApp,
	)

	sketchybarWindowID := getSketchybarWindowID(windowID)

	item.windowIDs[windowID] = workspaceID

	batches = batch(batches, s("--add", "item", sketchybarWindowID, position))
	batches = batch(batches, m(s("--set", sketchybarWindowID), windowItem.ToArgs()))

	return batches, sketchybarWindowID
}

func checker(batches Batches, position sketchybar.Position) (Batches, error) {
	updateEvent, err := args.BuildEvent()

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

	batches = batch(batches, s("--add", "item", aerospaceCheckerItemName, position))
	batches = batch(batches, m(s("--set", aerospaceCheckerItemName), checkerItem.ToArgs()))
	batches = batch(batches, s("--subscribe", aerospaceCheckerItemName,
		events.DisplayChange,
		events.SpaceWindowsChange,
		events.SystemWoke,
		events.FrontAppSwitched,
	))

	return batches, nil
}

type workspaceColors struct {
	backgroundColor string
	color           string
}

func (item AerospaceItem) getWorkspaceColors(isFocusedWorkspace bool) workspaceColors {
	backgroundColor := settings.Sketchybar.Aerospace.WorkspaceBackgroundColor
	color := settings.Sketchybar.Aerospace.WorkspaceColor

	if isFocusedWorkspace {
		backgroundColor = settings.Sketchybar.Aerospace.WorkspaceFocusedBackgroundColor
		color = settings.Sketchybar.Aerospace.WorkspaceFocusedColor
	}

	return workspaceColors{
		backgroundColor,
		color,
	}
}

type windowVisibility struct {
	width        *int
	show         string
	color        string
	focusedColor string
}

func (item AerospaceItem) getWindowVisibility(isFocusedWorkspace bool) *windowVisibility {
	width := pointer(32)
	show := "on"
	color := settings.Sketchybar.Aerospace.WindowColor
	focusedColor := settings.Sketchybar.Aerospace.WindowFocusedColor

	if !isFocusedWorkspace {
		width = pointer(0)
		show = "off"
		color = colors.Transparent
		focusedColor = colors.Transparent
	}

	return &windowVisibility{
		width,
		show,
		color,
		focusedColor,
	}
}

func (item AerospaceItem) handleWorkspaceChange(
	_ context.Context,
	batches Batches,
	_ string,
	focusedWorkspaceID string,
) Batches {
	tree := item.aerospace.GetTree()

	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			isFocusedWorkspace := workspace.Workspace == focusedWorkspaceID

			sketchybarWorkspaceID := getSketchybarWorkspaceID(workspace.Workspace)

			colors := item.getWorkspaceColors(isFocusedWorkspace)
			workspaceItem := &sketchybar.ItemOptions{
				Background: sketchybar.BackgroundOptions{
					Color: sketchybar.ColorOptions{
						Color: colors.backgroundColor,
					},
				},
				Icon: sketchybar.ItemIconOptions{
					Color: sketchybar.ColorOptions{
						Color: colors.color,
					},
				},
			}

			batches = batch(batches, m(
				s(
					"--animate",
					sketchybar.AnimationTanh,
					settings.Sketchybar.Aerospace.TransitionTime,
					"--set",
					sketchybarWorkspaceID,
				),
				workspaceItem.ToArgs(),
			))

			sketchybarBracketID := getSketchybarBracketID(workspace.Workspace)
			bracketItem := sketchybar.BracketOptions{
				Background: sketchybar.BackgroundOptions{
					Border: sketchybar.BorderOptions{
						Color: colors.backgroundColor,
					},
				},
			}

			batches = batch(batches, m(
				s(
					"--animate",
					sketchybar.AnimationTanh,
					settings.Sketchybar.Aerospace.TransitionTime,
					"--set",
					sketchybarBracketID,
				),
				bracketItem.ToArgs(),
			))

			for _, window := range workspace.Windows {
				sketchybarWindowID := getSketchybarWindowID(window)

				windowVisibility := item.getWindowVisibility(isFocusedWorkspace)
				windowItem := &sketchybar.ItemOptions{
					Display: strconv.Itoa(monitor.Monitor),
					Width:   windowVisibility.width,
					Icon: sketchybar.ItemIconOptions{
						Drawing: windowVisibility.show,
					},
				}

				batches = batch(batches, m(
					s(
						"--animate",
						sketchybar.AnimationTanh,
						settings.Sketchybar.Aerospace.TransitionTime,
						"--set",
						sketchybarWindowID,
					),
					windowItem.ToArgs(),
				))
			}
		}
	}
	return batches
}

func (item AerospaceItem) getSketchybarDisplayIndex(
	monitorCount int,
	monitorID aerospace.MonitorID,
) string {
	result := monitorID + 1

	// fmt.Printf("%d ==  %d\n", result, monitorCount)

	if result >= monitorCount {
		result = 1
	}

	return strconv.Itoa(result)
}

func (item AerospaceItem) handleDisplayChange(batches Batches) Batches {
	tree := item.aerospace.GetTree()

	for _, monitor := range tree.Monitors {
		for _, workspace := range monitor.Workspaces {
			sketchybarWorkspaceID := getSketchybarWorkspaceID(workspace.Workspace)

			workspaceItem := &sketchybar.ItemOptions{
				Display: item.getSketchybarDisplayIndex(len(tree.Monitors), monitor.Monitor),
			}

			batches = batch(batches, m(s("--set", sketchybarWorkspaceID), workspaceItem.ToArgs()))

			for _, windowID := range workspace.Windows {
				sketchybarWindowID := getSketchybarWindowID(windowID)

				windowItem := &sketchybar.ItemOptions{
					Display: item.getSketchybarDisplayIndex(len(tree.Monitors), monitor.Monitor),
				}

				batches = batch(batches, m(s("--set", sketchybarWindowID), windowItem.ToArgs()))
			}
		}
	}

	return batches
}

func (item *AerospaceItem) addWorkspaceBracket(
	batches Batches,
	workspaceID string,
	sketchybarWindowIDs []string,
) Batches {
	workspaceBracketItem := sketchybar.BracketOptions{
		Background: sketchybar.BackgroundOptions{
			Drawing: "on",
			Border: sketchybar.BorderOptions{
				Color: settings.Sketchybar.Aerospace.WorkspaceBackgroundColor,
			},
			Color: sketchybar.ColorOptions{
				Color: colors.Transparent,
			},
		},
	}

	sketchybarSpaceID := getSketchybarWorkspaceID(workspaceID)
	sketchybarBracketID := getSketchybarBracketID(workspaceID)

	batches = batch(batches, m(s(
		"--add",
		"bracket",
		sketchybarBracketID,
		sketchybarSpaceID),
		sketchybarWindowIDs,
	))

	batches = batch(batches, m(s(
		"--set",
		sketchybarBracketID,
	), workspaceBracketItem.ToArgs()))

	return batches
}

func (item *AerospaceItem) addWorkspaceSpacer(
	batches Batches,
	workspaceID string,
	position sketchybar.Position,
) Batches {
	workspaceSpacerItem := sketchybar.ItemOptions{
		Width: settings.Sketchybar.ItemSpacing,
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	sketchybarSpacerID := getSketchybarSpacerID(workspaceID)
	batches = batch(batches, s(
		"--add",
		"item",
		sketchybarSpacerID,
		position,
	))
	batches = batch(batches, m(s(
		"--set",
		sketchybarSpacerID,
	), workspaceSpacerItem.ToArgs()))

	return batches
}

func (item *AerospaceItem) applyTree(
	ctx context.Context,
	batches Batches,
	position sketchybar.Position,
) (Batches, error) {
	tree := item.aerospace.GetTree()
	focusedSpaceID := item.aerospace.GetFocusedWorkspaceID(ctx)

	aerospaceSpacerItem := sketchybar.ItemOptions{
		Width: settings.Sketchybar.ItemSpacing,
		Background: sketchybar.BackgroundOptions{
			Drawing: "off",
		},
	}

	sketchybarSpacerID := "aerospace.spacer"
	batches = batch(batches, s(
		"--add",
		"item",
		sketchybarSpacerID,
		position,
	))
	batches = batch(batches, m(s(
		"--set",
		sketchybarSpacerID,
	), aerospaceSpacerItem.ToArgs()))

	var aggregatedErr error
	for _, monitor := range tree.Monitors {
		item.logger.DebugContext(ctx, "monitor", slog.Int("monitor", monitor.Monitor))

		for _, workspace := range monitor.Workspaces {
			isFocusedWorkspace := focusedSpaceID == workspace.Workspace

			item.logger.DebugContext(
				ctx,
				"workspace",
				slog.String("workspace", workspace.Workspace),
				slog.Bool("focused", isFocusedWorkspace),
			)

			sketchybarSpaceID := getSketchybarWorkspaceID(workspace.Workspace)
			workspaceSpace, err := item.workspaceToSketchybar(
				isFocusedWorkspace,
				len(tree.Monitors),
				monitor.Monitor,
				workspace.Workspace,
			)

			if err != nil {
				aggregatedErr = errors.Join(aggregatedErr, err)
				continue
			}

			batches = batch(batches, s("--add", "item", sketchybarSpaceID, position))
			batches = batch(batches, m(s("--set", sketchybarSpaceID), workspaceSpace.ToArgs()))

			sketchybarWindowIDs := make([]string, len(workspace.Windows))
			for i, windowID := range workspace.Windows {
				window := tree.IndexedWindows[windowID]

				item.logger.DebugContext(
					ctx,
					"window",
					slog.Int("window", window.ID),
					slog.String("app", window.App),
				)

				var sketchybarWindowID string
				batches, sketchybarWindowID = item.addWindowToSketchybar(
					batches,
					position,
					isFocusedWorkspace,
					monitor.Monitor,
					workspace.Workspace,
					window.ID,
					window.App,
				)
				sketchybarWindowIDs[i] = sketchybarWindowID
			}

			batches = item.addWorkspaceBracket(batches, workspace.Workspace, sketchybarWindowIDs)
			batches = item.addWorkspaceSpacer(batches, workspace.Workspace, position)
		}
	}

	return batches, aggregatedErr
}

func (item *AerospaceItem) highlightWindows(
	batches Batches,
	workspaceID string,
	app string,
) Batches {
	windowsOfFocusedWorkspace := item.aerospace.WindowsOfWorkspace(workspaceID)

	for _, window := range windowsOfFocusedWorkspace {
		windowItemID := getSketchybarWindowID(window.ID)

		color := settings.Sketchybar.Aerospace.WindowColor
		if utils.Equals(window.App, app) {
			color = settings.Sketchybar.Aerospace.WindowFocusedColor
		}

		windowItem := sketchybar.ItemOptions{
			Icon: sketchybar.ItemIconOptions{
				Color: sketchybar.ColorOptions{
					Color: color,
				},
			},
		}

		batches = batch(batches, m(s(
			"--animate",
			sketchybar.AnimationTanh,
			settings.Sketchybar.Aerospace.TransitionTime,
			"--set",
			windowItemID,
		), windowItem.ToArgs()))
	}

	return batches
}

func isAerospace(name string) bool {
	return name == AerospaceName
}

var _ WentsketchyItem = (*AerospaceItem)(nil)
