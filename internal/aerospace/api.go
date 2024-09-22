package aerospace

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/utils"
)

type MonitorID = int
type WorkspaceID = string
type WindowID = int

type API interface {
	Monitors(ctx context.Context) ([]MonitorID, error)
	FocusedMonitor(ctx context.Context) (MonitorID, error)
	WorkspacesOfMonitor(ctx context.Context, monitorID int) ([]WorkspaceID, error)
	FocusedWorkspace(ctx context.Context) (WorkspaceID, error)
	WindowsOfWorkspace(ctx context.Context, workspaceID string) ([]*Window, error)
	WindowsOfMonitor(ctx context.Context, monitorID string) ([]*Window, error)
	FullWindows(ctx context.Context) ([]*FullWindow, error)
	FocusedWorkspaceWindows(ctx context.Context) ([]*Window, error)
	FocusedMonitorWindows(ctx context.Context) ([]*Window, error)
	FocusedWindow(ctx context.Context) (WindowID, error)
}

type realAPI struct {
	logger  *slog.Logger
	command *command.Command
}

func NewAPI(logger *slog.Logger, command *command.Command) API {
	return realAPI{
		logger,
		command,
	}
}

func (api realAPI) Monitors(ctx context.Context) ([]MonitorID, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-monitors",
		"--format",
		monitorOutputFormat(),
	)

	if err != nil {
		return make([]MonitorID, 0), fmt.Errorf("aerospace: could not get monitors. %w", err)
	}

	return splitAndMapMonitors(output)
}

func (api realAPI) FocusedMonitor(ctx context.Context) (MonitorID, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-monitors",
		"--focused",
		"--format",
		monitorOutputFormat(),
	)

	if err != nil {
		return 0, fmt.Errorf("aerospace: could not get monitors. %w", err)
	}

	monitors, err := splitAndMapMonitors(output)

	if err != nil {
		return 0, fmt.Errorf("aerospace: could not get split lines. %w", err)
	}

	if len(monitors) == 0 {
		return 0, fmt.Errorf("aerospace: could not find focused monitor. %w", err)
	}

	return monitors[0], nil
}

func (api realAPI) WorkspacesOfMonitor(ctx context.Context, monitorID int) ([]WorkspaceID, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-workspaces",
		"--monitor",
		strconv.Itoa(monitorID),
		"--format",
		workspaceOutputFormat(),
	)

	if err != nil {
		return make([]WorkspaceID, 0), fmt.Errorf("aerospace: could not get workspaces. %w", err)
	}

	return splitAndMapWorkspaces(output)
}

func (api realAPI) FocusedWorkspace(ctx context.Context) (WorkspaceID, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-workspaces",
		"--focused",
		"--format",
		workspaceOutputFormat(),
	)

	if err != nil {
		return "", fmt.Errorf("aerospace: could not get workspaces. %w", err)
	}

	workspaces, err := splitAndMapWorkspaces(output)

	if err != nil {
		return "", fmt.Errorf("aerospace: could not get split lines. %w", err)
	}

	if len(workspaces) == 0 {
		return "", fmt.Errorf("aerospace: could not find focused workspace. %w", err)
	}

	return workspaces[0], nil
}

func (api realAPI) WindowsOfWorkspace(ctx context.Context, workspaceID string) ([]*Window, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-windows",
		"--workspace",
		workspaceID,
		"--format",
		windowOutputFormat(),
	)

	if err != nil {
		return make([]*Window, 0), fmt.Errorf(
			"aerospace: could not get windows with workspace %s. %w",
			workspaceID,
			err,
		)
	}

	return splitAndMapWindows(output)
}

func (api realAPI) WindowsOfMonitor(ctx context.Context, monitorID string) ([]*Window, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-windows",
		"--monitor",
		monitorID,
		"--format",
		windowOutputFormat(),
	)

	if err != nil {
		return make([]*Window, 0), fmt.Errorf(
			"aerospace: could not get windows with workspace %s. %w",
			monitorID,
			err,
		)
	}

	return splitAndMapWindows(output)
}

func (api realAPI) FullWindows(ctx context.Context) ([]*FullWindow, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-windows",
		"--all",
		"--format",
		fullWindowOutputFormat(),
	)

	if err != nil {
		return make([]*FullWindow, 0), fmt.Errorf(
			"aerospace: could not get full windows. %w",
			err,
		)
	}

	return splitAndMapFullWindows(output)
}
func (api realAPI) FocusedWorkspaceWindows(ctx context.Context) ([]*Window, error) {
	return api.WindowsOfWorkspace(ctx, "focused")
}

func (api realAPI) FocusedMonitorWindows(ctx context.Context) ([]*Window, error) {
	return api.WindowsOfMonitor(ctx, "focused")
}

func (api realAPI) FocusedWindow(ctx context.Context) (WindowID, error) {
	output, err := api.command.Run(
		ctx,
		"aerospace",
		"list-windows",
		"--focused",
		"--format",
		outputFormatWindowID,
	)

	if err != nil {
		return 0, fmt.Errorf(
			"aerospace: could not focused windows. %w",
			err,
		)
	}

	windowIDs, err := splitAndMap(output, func(splitted []string) (WindowID, error) {
		id, err := strconv.Atoi(utils.Sanitize(splitted[0]))

		if err != nil {
			return 0, err
		}

		return id, nil
	})

	if err != nil {
		return 0, fmt.Errorf("aerospace: could not get split lines. %w", err)
	}

	if len(windowIDs) == 0 {
		return 0, fmt.Errorf("aerospace: could not find focused window. %w", err)
	}

	return windowIDs[0], nil
}

func splitAndMap[T any](output string, mapTo func([]string) (T, error)) ([]T, error) {
	lines := strings.Split(output, "\n")

	var aggregatedErr error
	var result = make([]T, len(lines)-1)
	for i, line := range lines {
		if line == "" {
			continue
		}

		mapped, err := mapTo(strings.Split(line, outputFormatSeparator))

		if err != nil {
			aggregatedErr = errors.Join(aggregatedErr, fmt.Errorf(
				"aerospace: could not parse line %s. %w",
				line,
				err,
			))

			continue
		}

		result[i] = mapped
	}

	return result, aggregatedErr
}

func splitAndMapWindows(output string) ([]*Window, error) {
	return splitAndMap(output, func(splitted []string) (*Window, error) {
		id, err := strconv.Atoi(utils.Sanitize(splitted[0]))

		if err != nil {
			return nil, err
		}

		return &Window{
			ID:  id,
			App: utils.Sanitize(splitted[1]),
		}, nil
	})
}

func splitAndMapFullWindows(output string) ([]*FullWindow, error) {
	return splitAndMap(output, func(splitted []string) (*FullWindow, error) {
		id, err := strconv.Atoi(utils.Sanitize(splitted[0]))

		if err != nil {
			return nil, err
		}

		monitorID, err := strconv.Atoi(utils.Sanitize(splitted[3]))

		if err != nil {
			return nil, err
		}

		return &FullWindow{
			ID:          id,
			App:         utils.Sanitize(splitted[1]),
			WorkspaceID: utils.Sanitize(splitted[2]),
			MonitorID:   monitorID,
		}, nil
	})
}

func splitAndMapMonitors(output string) ([]MonitorID, error) {
	return splitAndMap(output, func(splitted []string) (MonitorID, error) {
		id, err := strconv.Atoi(utils.Sanitize(splitted[0]))

		if err != nil {
			return 0, err
		}

		return id, nil
	})
}

func splitAndMapWorkspaces(output string) ([]WorkspaceID, error) {
	return splitAndMap(output, func(splitted []string) (WorkspaceID, error) {
		return utils.Sanitize(splitted[0]), nil
	})
}
