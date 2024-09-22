package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/aerospace"
	"github.com/lucax88x/wentsketchy/internal/aerospace/events"
	"github.com/lucax88x/wentsketchy/internal/fifo"
)

type FifoServer struct {
	logger    *slog.Logger
	config    *config.Config
	fifo      *fifo.Reader
	aerospace aerospace.Aerospace
}

func NewFifoServer(
	logger *slog.Logger,
	config *config.Config,
	fifo *fifo.Reader,
	aerospace aerospace.Aerospace,
) *FifoServer {
	return &FifoServer{
		logger,
		config,
		fifo,
		aerospace,
	}
}

func (f FifoServer) Start(ctx context.Context, fifoPath string) {
	ch := make(chan string)
	defer func() {
		close(ch)
	}()

	go func(ctx context.Context) {
		err := f.fifo.Listen(ctx, fifoPath, ch)

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not listen fifo", slog.Any("err", err))
		}
	}(ctx)

	for {
		select {
		case msg := <-ch:
			f.handle(
				ctx,
				msg,
			)
		case <-ctx.Done():
			f.logger.InfoContext(ctx, "server: cancel")
			return
		}
	}
}

func (f FifoServer) handle(
	ctx context.Context,
	msg string,
) {
	if strings.HasPrefix(msg, "init") {
		err := f.config.Init(ctx)

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not handle init", slog.Any("err", err))
		}

		return
	}

	if strings.HasPrefix(msg, "update") {
		args, err := args.FromEvent(msg)

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not get args", slog.Any("err", err))
		}

		f.logger.InfoContext(
			ctx,
			"server: react",
			slog.String("event", "update"),
			slog.String("sender", args.Name),
			slog.String("sender", args.Event),
			// slog.Any("args", args),
		)

		err = f.config.Update(ctx, args)

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not handle update", slog.Any("err", err))
		}

		return
	}

	if strings.HasPrefix(msg, string(events.WorkspaceChange)) {
		f.logger.InfoContext(
			ctx,
			"server: react",
			slog.String("event", string(events.WorkspaceChange)),
		)

		var data AerospaceWorkspaceChangeEvent

		eventJSON, _ := strings.CutPrefix(msg, string(events.WorkspaceChange))

		err := json.Unmarshal([]byte(eventJSON), &data)

		if err != nil {
			f.logger.ErrorContext(
				ctx,
				"server: could not deserialize data for aerospace_workspace_change",
				slog.String("msg", msg),
				slog.Any("err", err),
			)
		}

		f.aerospace.SetPrevWorkspaceID(data.Prev)
		f.aerospace.SetFocusedWorkspaceID(data.Focused)

		return
	}

	f.logger.InfoContext(ctx, "server: did not handle message", slog.String("msg", msg))
}

type AerospaceWorkspaceChangeEvent struct {
	Focused string `json:"focused"`
	Prev    string `json:"prev"`
}
