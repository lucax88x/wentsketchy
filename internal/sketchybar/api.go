package sketchybar

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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
	out, err := command.Run("sketchybar", flattenAndFix(arg...)...)

	api.logger.InfoContext(ctx, out)

	if err != nil {
		return fmt.Errorf("error while running sketchybar %w", err)
	}

	return nil
}

func flattenAndFix(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		for i, str := range slice {
			slice[i] = strings.TrimSpace(str)
			result = append(result, slice[i])
		}

	}
	return result
}
