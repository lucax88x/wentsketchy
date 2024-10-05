package aerospace

import (
	"context"
	"log/slog"
	"sort"
	"strconv"
)

type TreeBuilder interface {
	// Workspace id is ordered because of numbers, but windows cannot be ordered
	Build(ctx context.Context) (*Tree, error)
}

type realTreeBuilder struct {
	logger *slog.Logger
	api    API
}

func NewTreeBuilder(logger *slog.Logger, api API) TreeBuilder {
	return realTreeBuilder{
		logger,
		api,
	}
}

func (t realTreeBuilder) Build(ctx context.Context) (*Tree, error) {
	fullWorkspaces, err := t.api.FullWorkspaces(ctx)

	if err != nil {
		return nil, err
	}

	fullWindows, err := t.api.FullWindows(ctx)

	if err != nil {
		return nil, err
	}

	indexedMonitors := make(IndexedMonitors, 0)
	indexedWorkspaces := make(IndexedWorkspaces, 0)
	indexedWindows := make(IndexedWindows, 0)

	for _, fullWorkspace := range fullWorkspaces {
		_, foundWorkspace := indexedWorkspaces[fullWorkspace.ID]
		if !foundWorkspace {
			indexedWorkspaces[fullWorkspace.ID] = &WorkspaceWithWindowIDs{
				fullWorkspace.ID,
				make([]WindowID, 0),
			}
		}

		monitor, foundMonitor := indexedMonitors[fullWorkspace.MonitorID]
		if !foundMonitor {
			monitor = &MonitorWithWorkspaceIDs{
				fullWorkspace.MonitorID,
				make([]WorkspaceID, 0),
			}

			indexedMonitors[fullWorkspace.MonitorID] = monitor
		}

		if !containsString(monitor.Workspaces, fullWorkspace.ID) {
			monitor.Workspaces = append(monitor.Workspaces, fullWorkspace.ID)
		}
	}

	for _, fullWindow := range fullWindows {
		indexedWindows[fullWindow.ID] = &Window{
			ID:  fullWindow.ID,
			App: fullWindow.App,
		}

		workspace, foundWorkspace := indexedWorkspaces[fullWindow.WorkspaceID]
		if !foundWorkspace {
			t.logger.ErrorContext(
				ctx,
				"could not find workspace",
				slog.String("workspace", fullWindow.WorkspaceID),
				slog.Int("window", fullWindow.ID),
			)
		}

		workspace.Windows = append(workspace.Windows, fullWindow.ID)
	}

	sortWorkspaces(indexedMonitors)

	branches := make([]*Branch, 0)
	for _, monitor := range indexedMonitors {
		branchWorkspaces := make([]*WorkspaceWithWindowIDs, 0)

		for _, workspaceID := range monitor.Workspaces {
			workspace := indexedWorkspaces[workspaceID]

			branchWorkspaces = append(branchWorkspaces, &WorkspaceWithWindowIDs{
				Workspace: workspaceID,
				Windows:   workspace.Windows,
			})
		}

		branch := &Branch{
			Monitor:    monitor.Monitor,
			Workspaces: branchWorkspaces,
		}

		branches = append(branches, branch)
	}

	return &Tree{
		branches,
		indexedMonitors,
		indexedWorkspaces,
		indexedWindows,
	}, nil
}

func containsString(slice []string, e string) bool {
	for _, a := range slice {
		if a == e {
			return true
		}
	}
	return false
}

func indexWindows(windows []*Window) IndexedWindows {
	indexedWindows := make(IndexedWindows, len(windows))

	for _, window := range windows {
		indexedWindows[window.ID] = window
	}

	return indexedWindows
}

func indexFullWindows(windows []*FullWindow) IndexedFullWindows {
	indexedFullWindows := make(IndexedFullWindows, len(windows))

	for _, window := range windows {
		indexedFullWindows[window.ID] = window
	}

	return indexedFullWindows
}

func sortWorkspaces(indexedMonitors IndexedMonitors) {
	for _, monitor := range indexedMonitors {
		sort.Slice(monitor.Workspaces, func(i, j int) bool {
			left, _ := strconv.Atoi(monitor.Workspaces[i])
			right, _ := strconv.Atoi(monitor.Workspaces[j])

			return left < right
		})
	}
}

type IndexedMonitors = map[int]*MonitorWithWorkspaceIDs
type IndexedWorkspaces = map[string]*WorkspaceWithWindowIDs
type IndexedWindows = map[int]*Window
type IndexedFullWindows = map[int]*FullWindow

type Tree struct {
	Monitors []*Branch

	IndexedMonitors   IndexedMonitors
	IndexedWorkspaces IndexedWorkspaces
	IndexedWindows    IndexedWindows
}

type Branch struct {
	Monitor    MonitorID
	Workspaces []*WorkspaceWithWindowIDs
}

type MonitorWithWorkspaceIDs struct {
	Monitor    MonitorID
	Workspaces []WorkspaceID
}

type WorkspaceWithWindowIDs struct {
	Workspace WorkspaceID
	Windows   []WindowID
}
