package fifo

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"
)

type Reader struct {
	logger *slog.Logger
}

func NewFifoReader(logger *slog.Logger) *Reader {
	return &Reader{
		logger,
	}
}

func makeSureFifoExists(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		// TODO: check if existing is fifo!
		return nil
	}

	return syscall.Mkfifo(path, 0640)
}

func (f *Reader) Start(path string) error {
	err := makeSureFifoExists(path)

	if err != nil {
		return fmt.Errorf("fifo: error creating file. %w", err)
	}

	return nil
}

func (f *Reader) Listen(
	ctx context.Context,
	path string,
	ch chan<- string,
) error {
	pipe, err := os.OpenFile(path, os.O_RDWR, os.ModeNamedPipe)

	defer func() {
		err = pipe.Close()

		if err != nil {
			err = fmt.Errorf("fifo: could not close reader %w", err)
		}
	}()

	if err != nil {
		return fmt.Errorf("fifo: error opening for reading. %w", err)
	}

	reader := bufio.NewReader(pipe)

	internalCh := make(chan []byte)
	continueReading := true

	defer close(internalCh)

	go func() {
		for continueReading {
			line, readErr := reader.ReadBytes('\n')

			if readErr != nil {
				err = readErr
				f.logger.ErrorContext(ctx, "fifo: readbytes err", slog.Any("error", err))
				break
			}

			if continueReading {
				internalCh <- line
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			f.logger.InfoContext(ctx, "fifo: cancel")
			continueReading = false

			err = pipe.Close()

			if err != nil {
				err = fmt.Errorf("fifo: could not close reader %w", err)
			}
			return nil
		case data := <-internalCh:
			nline := string(data)
			nline = strings.TrimRight(nline, "\r\n")
			ch <- nline
		}
	}
}
