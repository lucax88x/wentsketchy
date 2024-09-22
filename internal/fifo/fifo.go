package fifo

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"syscall"
)

const Separator = '¬'

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
		err = os.Remove(path)

		if err != nil {
			return fmt.Errorf("fifo: could not remove fifo. %w", err)
		}
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
			line, readErr := reader.ReadBytes(Separator)

			if errors.Is(err, io.EOF) {
				f.logger.ErrorContext(ctx, "fifo: got EOF")
				break
			}

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

			err = ensureClose(path)

			if err != nil {
				err = fmt.Errorf("fifo: could not close fifo with EOF %w", err)
			}

			return nil
		case data := <-internalCh:
			nline := string(data)
			nline = strings.TrimRight(nline, string(Separator))
			nline = strings.TrimLeft(nline, "\n")
			ch <- nline
		}
	}
}

func ensureClose(path string) error {
	// f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	//
	// if err != nil {
	// 	return fmt.Errorf("fifo: error opening file to write EOF: %w", err)
	// }
	//
	// _, err = f.WriteString("EOF")
	//
	// if err != nil {
	// 	return fmt.Errorf("fifo: error while writing EOF: %w", err)
	// }

	err := os.Remove(path)

	if err != nil {
		return fmt.Errorf("fifo: could not remove fifo. %w", err)
	}

	return nil
}
