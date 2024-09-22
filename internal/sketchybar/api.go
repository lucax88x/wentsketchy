package sketchybar

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lucax88x/wentsketchy/internal/command"
	"github.com/lucax88x/wentsketchy/internal/sketchybar/query"
)

type API interface {
	QueryBar(ctx context.Context) (query.Bar, error)
	Run(ctx context.Context, arg []string) error
}

type realAPI struct {
	logger  *slog.Logger
	command *command.Command
}

func NewAPI(logger *slog.Logger, command *command.Command) realAPI {
	return realAPI{
		logger,
		command,
	}
}

func (api realAPI) Run(ctx context.Context, arg []string) error {
	_, err := api.run(ctx, arg...)

	return err
}

func (api realAPI) QueryBar(ctx context.Context) (query.Bar, error) {
	var bar query.Bar

	out, err := api.run(ctx, "--query", "bar")

	if err != nil {
		return bar, fmt.Errorf("sketchybar: cannot deserialize query bar. %w", err)
	}

	err = json.Unmarshal([]byte(out), &bar)

	if err != nil {
		return bar, fmt.Errorf("sketchybar: cannot deserialize query bar. %w", err)
	}

	return bar, nil
}

func (api realAPI) run(ctx context.Context, arg ...string) (string, error) {
	if len(arg) == 0 {
		return "", nil
	}

	out, err := api.command.Run(ctx, "sketchybar", flattenAndFix(arg)...)

	if err != nil {
		api.logger.ErrorContext(ctx, out)
		return "", fmt.Errorf("sketchybar: error while running. %w", err)
	}

	// api.logger.InfoContext(ctx, out)

	return out, nil
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
