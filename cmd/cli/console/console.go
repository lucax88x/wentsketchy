package console

import (
	"io"
)

type Console struct {
	Stdout io.Writer
	Stderr io.Writer
}
