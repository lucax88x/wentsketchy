package aerospace

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/lucax88x/wentsketchy/internal/command"
)

type API interface {
	Monitors(ctx context.Context) ([]*Monitor, error)
	Workspaces(_ context.Context, monitorID int) ([]*Workspace, error)
}

type realAPI struct {
	logger *slog.Logger
}

func NewAPI(logger *slog.Logger) API {
	return realAPI{
		logger,
	}
}

func (api realAPI) Monitors(_ context.Context) ([]*Monitor, error) {
	output, err := command.Run("aerospace", "list-monitors")

	if err != nil {
		return make([]*Monitor, 0), fmt.Errorf("could not get monitors from aerospace %w", err)
	}

	return splitAndMap(output, func(splitted []string) (*Monitor, error) {
		id, err := strconv.Atoi(strings.TrimSpace(splitted[0]))

		if err != nil {
			return nil, err
		}

		return &Monitor{
			Id:   id,
			Name: strings.TrimSpace(splitted[1]),
		}, nil
	})
}

func (api realAPI) Workspaces(_ context.Context, monitorID int) ([]*Workspace, error) {
	output, err := command.Run("aerospace", "list-workspaces", "--monitor", strconv.Itoa(monitorID))

	if err != nil {
		return make([]*Workspace, 0), fmt.Errorf("could not get workspaces from aerospace %w", err)
	}

	return splitAndMap(output, func(splitted []string) (*Workspace, error) {
		return &Workspace{
			Id: splitted[0],
		}, nil
	})
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
				"could not parse line %s for %w",
				line,
				err,
			))

			continue
		}

		result[i] = mapped
	}

	return result, aggregatedErr
}