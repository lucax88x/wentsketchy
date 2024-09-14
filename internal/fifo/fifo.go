package fifo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
)

func makeSureFifoExists(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		return nil
	}

	return syscall.Mkfifo(path, 0640)
}

func Start(path string) error {
	err := makeSureFifoExists(path)

	if err != nil {
		return fmt.Errorf("fifo: error creating FIFO. %w", err)
	}

	return nil
}

func Listen(
	path string,
	ch chan<- string,
	wg *sync.WaitGroup,
) error {
	defer func() {
		wg.Done()
	}()

	pipe, err := os.OpenFile(path, os.O_RDWR, 0640)

	defer func() {
		err = pipe.Close()

		err = fmt.Errorf("fifo: could not close reader %w", err)
	}()

	if err != nil {
		return fmt.Errorf("fifo: error opening FIFO for reading. %w", err)
	}

	reader := bufio.NewReader(pipe)

	for {
		line, readErr := reader.ReadBytes('\n')

		if readErr != nil {
			err = readErr
			break
		}

		// Remove new line char
		nline := string(line)
		nline = strings.TrimRight(nline, "\r\n")

		ch <- nline
	}

	if err != nil {
		return fmt.Errorf("error opening FIFO for reading. %w", err)
	}

	return nil
}
