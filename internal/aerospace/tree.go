package aerospace

import (
	"context"
	"errors"
	"log/slog"
	"sort"
	"sync"
)

type Tree interface {
	// Workspace id is ordered because of numbers, but windows cannot be ordered
	Build(ctx context.Context) (*AerospaceTree, error)
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
	ctx context.Context,
	monitor *Monitor,
	ch chan<- *WorkspacesResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	workspaces, err := t.api.Workspaces(ctx, monitor.Id)

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
	ctx context.Context,
	workspace *Workspace,
	ch chan<- *WindowsResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	windows, err := t.api.Windows(ctx, workspace.Id)

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
	ctx context.Context,
	workspaces []*Workspace,
) ([]*WorkspaceWithWindows, error) {
	t.logger.InfoContext(ctx, "fetching for workspaces", slog.Int("workspaces", len(workspaces)))

	var wg sync.WaitGroup
	windowsCh := make(chan *WindowsResult, len(workspaces))
	wg.Add(len(workspaces))

	for _, workspace := range workspaces {
		go t.fetchWindows(ctx, workspace, windowsCh, &wg)
	}

	go func() {
		wg.Wait()
		close(windowsCh)
	}()

	var result []*WorkspaceWithWindows
	var aggregatedErr error
	for message := range windowsCh {
		if message.Error != nil {
			aggregatedErr = errors.Join(aggregatedErr, message.Error)
		}

		result = append(result, &WorkspaceWithWindows{
			Workspace: message.Workspace,
			Windows:   message.Windows,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Workspace.Id < result[j].Workspace.Id
	})

	return result, aggregatedErr
}

func (t realTree) fetchMonitorsWorkspaces(
	ctx context.Context,
	monitors []*Monitor,
) ([]*MonitorWithWorkspaces, error) {
	t.logger.InfoContext(ctx, "fetching for monitors", slog.Int("monitors", len(monitors)))
	var wg sync.WaitGroup
	workspacesCh := make(chan *WorkspacesResult, len(monitors))
	wg.Add(len(monitors))

	for _, monitor := range monitors {
		go t.fetchWorkspaces(ctx, monitor, workspacesCh, &wg)
	}

	go func() {
		wg.Wait()
		close(workspacesCh)
	}()

	var result []*MonitorWithWorkspaces
	var aggregatedErr error
	for message := range workspacesCh {
		if message.Error != nil {
			aggregatedErr = errors.Join(aggregatedErr, message.Error)
		}

		workspacesWithWindows, err := t.fetchWorkspacesWindows(ctx, message.Workspaces)

		if err != nil {
			aggregatedErr = errors.Join(aggregatedErr, err)
		}

		result = append(result, &MonitorWithWorkspaces{
			Monitor:    message.Monitor,
			Workspaces: workspacesWithWindows,
		})
	}

	return result, aggregatedErr
}

func (t realTree) Build(ctx context.Context) (*AerospaceTree, error) {
	monitors, err := t.api.Monitors(ctx)

	if err != nil {
		return nil, err
	}

	monitorsWithWorkspaces, err := t.fetchMonitorsWorkspaces(ctx, monitors)

	if err != nil {
		return nil, err
	}

	return &AerospaceTree{
		Monitors: monitorsWithWorkspaces,
	}, nil
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

type WorkspaceWithWindows struct {
	Workspace *Workspace
	Windows   []*Window
}

type AerospaceTree struct {
	Monitors []*MonitorWithWorkspaces
}

type MonitorWithWorkspaces struct {
	Monitor    *Monitor
	Workspaces []*WorkspaceWithWindows
}
