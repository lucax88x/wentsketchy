package sketchybar

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/internal/command"
)

type API interface {
	Run(ctx context.Context, arg ...[]string) error
}

type realAPI struct {
	logger *slog.Logger
}

func NewAPI(logger *slog.Logger) API {
	return realAPI{
		logger,
	}
}

func (api realAPI) Run(ctx context.Context, arg ...[]string) error {
	a := flatten(arg...)

	out, err := command.Run("sketchybar", a...)

	api.logger.InfoContext(ctx, out)

	if err != nil {
		return fmt.Errorf("error while running sketchybar %w", err)
	}

	return nil
}

func flatten(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
