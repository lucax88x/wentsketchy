package server

import (
	"context"
	"log/slog"
	"strings"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/fifo"
)

type FifoServer struct {
	logger *slog.Logger
	config *config.Config
	fifo   *fifo.Reader
}

func NewFifoServer(
	logger *slog.Logger,
	config *config.Config,
	fifo *fifo.Reader,
) *FifoServer {
	return &FifoServer{
		logger,
		config,
		fifo,
	}
}

func (f FifoServer) Start(ctx context.Context, fifoPath string) {
	ch := make(chan string)
	defer func() {
		close(ch)
	}()

	f.config.SetFifoPath(fifoPath)

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
		err := f.config.Update(ctx, args.FromMsg(msg))

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not handle init", slog.Any("err", err))
		}

		return
	}

	if strings.HasPrefix(msg, "aerospace_workspace_change") {
		f.logger.InfoContext(ctx, "server: handling but only reading message", slog.String("msg", msg))

		return
	}

	f.logger.InfoContext(ctx, "server: did not handle message", slog.String("msg", msg))
}
