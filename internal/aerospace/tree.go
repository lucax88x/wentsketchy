package aerospace

import (
	"errors"
	"log/slog"
	"sort"
	"sync"
)

type Tree interface {
	// Workspace id is ordered because of numbers, but windows cannot be ordered
	Build() (*AerospaceTree, error)
}

type realTree struct {
	logger *slog.Logger
	api    API
}

func NewTree(logger *slog.Logger, api API) Tree {
	return realTree{
		logger,
		api,
	}
}

func (t realTree) fetchWorkspaces(
	monitor *Monitor,
	ch chan<- *WorkspacesResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	workspaces, err := t.api.Workspaces(monitor.ID)

	if err != nil {
		ch <- &WorkspacesResult{
			Workspaces: make([]*Workspace, 0),
			Error:      nil,
		}
		return
	}

	ch <- &WorkspacesResult{
		Monitor:    monitor,
		Workspaces: workspaces,
		Error:      nil,
	}
}

func (t realTree) fetchWindows(
	workspace *Workspace,
	ch chan<- *WindowsResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	windows, err := t.api.Windows(workspace.ID)

	if err != nil {
		ch <- &WindowsResult{
			Workspace: workspace,
			Windows:   make([]*Window, 0),
			Error:     nil,
		}
		return
	}

	ch <- &WindowsResult{
		Workspace: workspace,
		Windows:   windows,
		Error:     nil,
	}
}

func (t realTree) fetchWorkspacesWindows(
	workspaces []*Workspace,
	indexedWorkspaces IndexedWorkspaces,
	indexedWindows IndexedWindows,
) ([]*WorkspaceIDWithWindowIDs, error) {
	var wg sync.WaitGroup
	windowsCh := make(chan *WindowsResult, len(workspaces))
	wg.Add(len(workspaces))

	for _, workspace := range workspaces {
		go t.fetchWindows(workspace, windowsCh, &wg)
	}

	go func() {
		wg.Wait()
		close(windowsCh)
	}()

	var result []*WorkspaceIDWithWindowIDs
	var aggregatedErr error
	for message := range windowsCh {
		if message.Error != nil {
			aggregatedErr = errors.Join(aggregatedErr, message.Error)
		}

		windowIDs := make([]int, len(message.Windows))
		for i, window := range message.Windows {
			indexedWindows[window.ID] = window
			windowIDs[i] = window.ID
		}

		indexedWorkspaces[message.Workspace.ID] = &WorkspaceWithWindowIDs{
			message.Workspace,
			windowIDs,
		}

		result = append(result, &WorkspaceIDWithWindowIDs{
			Workspace: message.Workspace.ID,
			Windows:   windowIDs,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Workspace < result[j].Workspace
	})

	return result, aggregatedErr
}

func (t realTree) fetchMonitorsWorkspaces(
	monitors []*Monitor,
	indexedMonitors IndexedMonitors,
	indexedWorkspaces IndexedWorkspaces,
	indexedWindows IndexedWindows,
) ([]*MonitorIDWithWorkspaceIDs, error) {
	var wg sync.WaitGroup
	workspacesCh := make(chan *WorkspacesResult, len(monitors))
	wg.Add(len(monitors))

	for _, monitor := range monitors {
		go t.fetchWorkspaces(
			monitor,
			workspacesCh,
			&wg,
		)
	}

	go func() {
		wg.Wait()
		close(workspacesCh)
	}()

	var result []*MonitorIDWithWorkspaceIDs
	var aggregatedErr error
	for message := range workspacesCh {
		if message.Error != nil {
			aggregatedErr = errors.Join(aggregatedErr, message.Error)
		}

		workspacesWithWindows, err := t.fetchWorkspacesWindows(
			message.Workspaces,
			indexedWorkspaces,
			indexedWindows,
		)

		if err != nil {
			aggregatedErr = errors.Join(aggregatedErr, err)
		}

		workspaceIDs := make([]string, len(workspacesWithWindows))
		for i, workspacesWithWindows := range workspacesWithWindows {
			workspaceIDs[i] = workspacesWithWindows.Workspace
		}

		indexedMonitors[message.Monitor.ID] = &MonitorWithWorkspaceIDs{
			message.Monitor,
			workspaceIDs,
		}

		result = append(result, &MonitorIDWithWorkspaceIDs{
			Monitor:    message.Monitor.ID,
			Workspaces: workspacesWithWindows,
		})
	}

	return result, aggregatedErr
}

func (t realTree) Build() (*AerospaceTree, error) {
	monitors, err := t.api.Monitors()

	if err != nil {
		return nil, err
	}

	indexedMonitors := make(IndexedMonitors, 0)
	indexedWorkspaces := make(IndexedWorkspaces, 0)
	indexedWindows := make(IndexedWindows, 0)

	monitorsWithWorkspaces, err := t.fetchMonitorsWorkspaces(
		monitors,
		indexedMonitors,
		indexedWorkspaces,
		indexedWindows,
	)

	if err != nil {
		return nil, err
	}

	return &AerospaceTree{
		monitorsWithWorkspaces,
		indexedMonitors,
		indexedWorkspaces,
		indexedWindows,
	}, nil
}

func (t realTree) WindowsOfFocusedWorkspace() ([]*Window, error) {
	focusedWorkspace, err := t.api.FocusedWorkspace()

	if err != nil {
		return nil, err
	}

	return t.api.Windows(focusedWorkspace.ID)
}

type WorkspacesResult struct {
	Monitor    *Monitor
	Workspaces []*Workspace
	Error      error
}

type WindowsResult struct {
	Workspace *Workspace
	Windows   []*Window
	Error     error
}

type IndexedMonitors = map[int]*MonitorWithWorkspaceIDs
type IndexedWorkspaces = map[string]*WorkspaceWithWindowIDs
type IndexedWindows = map[int]*Window

type AerospaceTree struct {
	Monitors []*MonitorIDWithWorkspaceIDs

	IndexedMonitors   IndexedMonitors
	IndexedWorkspaces IndexedWorkspaces
	IndexedWindows    IndexedWindows
}

type Data struct {
	FocusedWorkspaceID string
	Tree               *AerospaceTree
}

type MonitorIDWithWorkspaceIDs struct {
	Monitor    int
	Workspaces []*WorkspaceIDWithWindowIDs
}

type WorkspaceIDWithWindowIDs struct {
	Workspace string
	Windows   []int
}

type MonitorWithWorkspaceIDs struct {
	Monitor    *Monitor
	Workspaces []string
}

type WorkspaceWithWindowIDs struct {
	Workspace *Workspace
	Windows   []int
}

func (data *Data) WindowsOfFocusedWorkspace(workspaceID string) []*Window {
	workspace, found := data.Tree.IndexedWorkspaces[workspaceID]
	if !found {
		return make([]*Window, 0)
	}

	windows := make([]*Window, 0, len(workspace.Windows))
	for _, windowID := range workspace.Windows {
		window, foundWindow := data.Tree.IndexedWindows[windowID]

		if !foundWindow {
			// log
			continue
		}

		windows = append(windows, window)
	}
	return windows
}
