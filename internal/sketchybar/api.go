package sketchybar

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lucax88x/wentsketchy/internal/command"
)

type API interface {
	Run(ctx context.Context, arg []string) error
}

type realAPI struct {
	logger *slog.Logger
}

func NewAPI(logger *slog.Logger) realAPI {
	return realAPI{
		logger,
	}
}

func (api realAPI) Run(ctx context.Context, arg []string) error {
	if len(arg) == 0 {
		api.logger.InfoContext(ctx, "sketchybar: there were no arguments, skipping")
		return nil
	}

	out, err := command.Run("sketchybar", flattenAndFix(arg)...)

	if err != nil {
		return fmt.Errorf("sketchybar: error while running. %w", err)
	}

	api.logger.InfoContext(ctx, out)

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
