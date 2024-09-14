package server

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/lucax88x/wentsketchy/cmd/cli/config"
	"github.com/lucax88x/wentsketchy/cmd/cli/config/args"
	"github.com/lucax88x/wentsketchy/internal/fifo"
)

type FifoServer struct {
	logger *slog.Logger
	config *config.Config
}

func NewFifoServer(logger *slog.Logger, config *config.Config) FifoServer {
	return FifoServer{
		logger,
		config,
	}
}

func (f FifoServer) Start(ctx context.Context, fifoPath string) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan string)
	defer func() {
		close(ch)
	}()

	f.config.SetFifoPath(fifoPath)

	go func() {
		err := fifo.Listen(fifoPath, ch, &wg)

		if err != nil {
			f.logger.ErrorContext(ctx, "server: could not listen fifo", slog.Any("err", err))
		}
	}()

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

	f.logger.InfoContext(ctx, "server: did not handle message", slog.String("msg", msg))
}
