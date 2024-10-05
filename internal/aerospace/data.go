package aerospace

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/singleflight"
)

type Aerospace interface {
	GetTree() *Tree
	GetPrevWorkspaceID() string
	SetPrevWorkspaceID(workspaceID string)
	GetFocusedWorkspaceID(ctx context.Context) string
	SetFocusedWorkspaceID(workspaceID string)
	GetFocusedMonitorID(ctx context.Context) int
	SetFocusedMonitorID(monitorID int)
	GetFocusedApp() string
	SetFocusedApp(app string)

	SingleFlightRefreshTree()

	FocusedMonitor(ctx context.Context) (MonitorID, error)
	WindowsOfWorkspace(workspaceID string) []*Window
	WindowsOfFocusedWorkspace(ctx context.Context) (IndexedWindows, error)
	WindowsOfFocusedMonitor(ctx context.Context) (IndexedWindows, error)
	FocusedWindow(ctx context.Context) (WindowID, error)
	AllFullWindows(ctx context.Context) (IndexedFullWindows, error)
}

type Data struct {
	logger      *slog.Logger
	api         API
	treeBuilder TreeBuilder

	prevWorkspaceID    string
	focusedWorkspaceID string
	prevMonitorID      int
	focusedMonitorID   int
	focusedApp         string
	tree               *Tree

	refreshTree *singleflight.Group
}

func New(
	logger *slog.Logger,
	api API,
	treeBuilder TreeBuilder,
) *Data {
	var g singleflight.Group
	return &Data{
		logger:      logger,
		api:         api,
		treeBuilder: treeBuilder,
		refreshTree: &g,
	}
}

func (data *Data) SingleFlightRefreshTree() {
	data.logger.Info("aerospace: refreshing..")
	_, err, shared := data.refreshTree.Do("refresh-aerospace-tree", data.refreshAerospaceData)
	data.logger.Info("aerospace: refreshed", slog.Bool("shared", shared))

	if err != nil {
		data.logger.Error("aerospace: error while refreshing tree", slog.Any("err", err))
	}
}

func (data *Data) GetTree() *Tree {
	return data.tree
}

func (data *Data) GetPrevWorkspaceID() string {
	return data.prevWorkspaceID
}

func (data *Data) GetFocusedWorkspaceID(ctx context.Context) string {
	if data.focusedWorkspaceID == "" {
		data.logger.InfoContext(ctx, "aerospace: no focused workspace, getting from aerospace")
		focusedWorkspaceID, err := data.api.FocusedWorkspace(ctx)

		if err != nil {
			data.logger.ErrorContext(ctx, "aerospace: could not get focused workspace")
			return ""
		}

		data.SetFocusedWorkspaceID(focusedWorkspaceID)
	}

	return data.focusedWorkspaceID
}

func (data *Data) SetPrevWorkspaceID(workspaceID string) {
	data.prevWorkspaceID = workspaceID
}

func (data *Data) SetFocusedWorkspaceID(workspaceID string) {
	data.focusedWorkspaceID = workspaceID
}

func (data *Data) SetFocusedMonitorID(monitorID int) {
	if data.focusedMonitorID != 0 {
		data.prevMonitorID = data.focusedMonitorID
	}
	data.focusedMonitorID = monitorID
}

func (data *Data) GetFocusedMonitorID(ctx context.Context) int {
	if data.focusedMonitorID == 0 {
		data.logger.InfoContext(ctx, "aerospace: no focused monitor, getting from aerospace")
		focusedMonitorID, err := data.api.FocusedMonitor(ctx)

		if err != nil {
			data.logger.ErrorContext(ctx, "aerospace: could not get focused monitor")
			return 0
		}

		data.SetFocusedMonitorID(focusedMonitorID)
	}

	return data.focusedMonitorID
}

func (data *Data) SetFocusedApp(app string) {
	data.focusedApp = app
}

func (data *Data) GetFocusedApp() string {
	return data.focusedApp
}

func (data *Data) FocusedMonitor(ctx context.Context) (MonitorID, error) {
	monitorID, err := data.api.FocusedMonitor(ctx)

	if err != nil {
		return 0, fmt.Errorf("aerospace: could not get focused monitor. %w", err)
	}

	return monitorID, nil
}

func (data *Data) WindowsOfFocusedWorkspace(ctx context.Context) (IndexedWindows, error) {
	windows, err := data.api.FocusedWorkspaceWindows(ctx)

	if err != nil {
		return make(IndexedWindows, 0), fmt.Errorf("aerospace: could not get focused windows. %w", err)
	}

	return indexWindows(windows), nil
}

func (data *Data) WindowsOfFocusedMonitor(ctx context.Context) (IndexedWindows, error) {
	windows, err := data.api.FocusedMonitorWindows(ctx)

	if err != nil {
		return make(IndexedWindows, 0), fmt.Errorf("aerospace: could not get focused windows. %w", err)
	}

	return indexWindows(windows), nil
}

func (data *Data) AllFullWindows(ctx context.Context) (IndexedFullWindows, error) {
	windows, err := data.api.FullWindows(ctx)

	if err != nil {
		return make(IndexedFullWindows, 0), fmt.Errorf("aerospace: could not get all windows. %w", err)
	}

	return indexFullWindows(windows), nil
}

func (data *Data) WindowsOfWorkspace(workspaceID string) []*Window {
	workspace, found := data.tree.IndexedWorkspaces[workspaceID]
	if !found {
		return make([]*Window, 0)
	}

	windows := make([]*Window, 0, len(workspace.Windows))
	for _, windowID := range workspace.Windows {
		window, foundWindow := data.tree.IndexedWindows[windowID]

		if !foundWindow {
			// log
			continue
		}

		windows = append(windows, window)
	}
	return windows
}

func (data *Data) FocusedWindow(ctx context.Context) (WindowID, error) {
	windowID, err := data.api.FocusedWindow(ctx)

	if err != nil {
		return 0, fmt.Errorf("aerospace: could not get focused windows. %w", err)
	}

	return windowID, nil
}

func (data *Data) refreshAerospaceData() (interface{}, error) {
	ctx := context.Background()

	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		data.logger.Info("aerospace: refresh took", slog.Duration("elapsed", elapsed))
	}()

	tree, err := data.treeBuilder.Build(ctx)

	if err != nil {
		return false, fmt.Errorf("aerospace: could not refresh tree. %w", err)
	}

	data.tree = tree
	return true, nil
}
