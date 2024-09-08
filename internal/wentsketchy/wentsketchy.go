package wentsketchy

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lucax88x/wentsketchy/internal/clock"
)

type Wentsketchy struct {
	Logger *slog.Logger
	Clock  clock.Clock
}

func NewWentsketchy(
	ctx context.Context,
	logger *slog.Logger,
	clock clock.Clock,
) (*Wentsketchy, error) {
	di := &Wentsketchy{
		Logger: logger,
		Clock:  clock,
	}

	err := initialize(ctx, di)

	if err != nil {
		return nil, fmt.Errorf("init: could not initialize wentsketchy. %w", err)
	}

	return di, nil
}

func initialize(_ context.Context, _ *Wentsketchy) error {
	return fmt.Errorf("TODO")
}
