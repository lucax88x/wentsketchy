package aerospace

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/internal/command"
)

type API interface {
	Monitors() ([]*Monitor, error)
	Workspaces(monitorID int) ([]*Workspace, error)
	FocusedWorkspace() (*Workspace, error)
	Windows(workspaceID string) ([]*Window, error)
}

type realAPI struct {
	logger *slog.Logger
}

func NewAPI(logger *slog.Logger) API {
	return realAPI{
		logger,
	}
}

func (api realAPI) Monitors() ([]*Monitor, error) {
	output, err := command.Run("aerospace", "list-monitors")

	if err != nil {
		return make([]*Monitor, 0), fmt.Errorf("aerospace: could not get monitors. %w", err)
	}

	return splitAndMap(output, func(splitted []string) (*Monitor, error) {
		id, err := strconv.Atoi(sanitize(splitted[0]))

		if err != nil {
			return nil, err
		}

		return &Monitor{
			ID:   id,
			Name: sanitize(splitted[1]),
		}, nil
	})
}

func (api realAPI) Workspaces(monitorID int) ([]*Workspace, error) {
	output, err := command.Run("aerospace", "list-workspaces", "--monitor", strconv.Itoa(monitorID))

	if err != nil {
		return make([]*Workspace, 0), fmt.Errorf("aerospace: could not get workspaces. %w", err)
	}

	return splitAndMap(output, func(splitted []string) (*Workspace, error) {
		return &Workspace{
			ID: sanitize(splitted[0]),
		}, nil
	})
}

func (api realAPI) FocusedWorkspace() (*Workspace, error) {
	output, err := command.Run("aerospace", "list-workspaces", "--focused")

	if err != nil {
		return nil, fmt.Errorf("aerospace: could not get workspaces. %w", err)
	}

	splitted, err := splitAndMap(output, func(splitted []string) (*Workspace, error) {
		return &Workspace{
			ID: sanitize(splitted[0]),
		}, nil
	})

	if err != nil {
		return nil, fmt.Errorf("aerospace: could not get split lines. %w", err)
	}

	if len(splitted) == 0 {
		return nil, fmt.Errorf("aerospace: could not find focused workspace. %w", err)
	}

	return splitted[0], nil
}

func (api realAPI) Windows(workspaceID string) ([]*Window, error) {
	output, err := command.Run("aerospace", "list-windows", "--workspace", workspaceID)

	if err != nil {
		return make([]*Window, 0), fmt.Errorf(
			"aerospace: could not get windows with workspace %s. %w",
			workspaceID,
			err,
		)
	}

	return splitAndMapWindows(output)
}

func splitAndMap[T any](output string, mapTo func([]string) (T, error)) ([]T, error) {
	lines := strings.Split(output, "\n")

	var aggregatedErr error
	var result = make([]T, len(lines)-1)
	for i, line := range lines {
		if line == "" {
			continue
		}

		mapped, err := mapTo(strings.Split(line, "|"))

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

func sanitize(str string) string {
	return strings.TrimSpace(str)
}

func splitAndMapWindows(output string) ([]*Window, error) {
	return splitAndMap(output, func(splitted []string) (*Window, error) {
		id, err := strconv.Atoi(sanitize(splitted[0]))

		if err != nil {
			return nil, err
		}

		return &Window{
			ID:    id,
			App:   sanitize(splitted[1]),
			Title: sanitize(splitted[2]),
		}, nil
	})
}
