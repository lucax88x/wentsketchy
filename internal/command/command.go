package command

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

type Command struct {
	logger *slog.Logger
}

func NewCommand(logger *slog.Logger) *Command {
	return &Command{
		logger,
	}
}

func (c Command) Run(ctx context.Context, name string, arg ...string) (string, error) {
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		c.logger.InfoContext(ctx, "command: took", slog.String("name", name), slog.Duration("elapsed", elapsed))
	}()

	cmd := exec.Command(name, arg...)

	// fmt.Println(fmt.Sprintf("%s %v", cmd.Path, cmd.Args))

	out, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("could not run command %w", err)
	}

	return string(out), nil
}
