package sketchybar

import "fmt"

type PaddingOptions struct {
	Left  int
	Right int
}

type ColorOptions struct {
	Color string
}

type BorderOptions struct {
	Width int
	Color string
}

type FontOptions struct {
	Font string
	Kind string
	Size string
}

var EmptyFontOptions = FontOptions{}

func (opts FontOptions) String() string {
	return fmt.Sprintf("%s:%s:%s", opts.Font, opts.Kind, opts.Size)
}

func (opts PaddingOptions) ToArgs(parent *string) []string {
	args := []string{}

	if opts.Right != 0 {
		args = withParent(args, parent, "padding_right=%d", opts.Right)
	}
	if opts.Left != 0 {
		args = withParent(args, parent, "padding_left=%d", opts.Left)
	}

	return args
}

func (opts ColorOptions) ToArgs(parent *string) []string {
	args := []string{}

	if opts.Color != "" {
		args = withParent(args, parent, "color=%s", opts.Color)
	}

	return args
}

func (opts BorderOptions) ToArgs(parent *string) []string {
	args := []string{}

	if opts.Width != 0 {
		args = withParent(args, parent, "border_width=%d", opts.Width)
	}
	if opts.Color != "" {
		args = withParent(args, parent, "border_color=%s", opts.Color)
	}

	return args
}

func withParent[T any](args []string, parent *string, format string, value T) []string {
	if parent != nil {
		format = *parent + "." + format
	}

	return append(args, fmt.Sprintf(format, value))
}

func with[T any](args []string, format string, value T) []string {
	return append(args, fmt.Sprintf(format, value))
}
