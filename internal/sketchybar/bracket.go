package sketchybar

type BracketOptions struct {
	// PADDING NOT SUPPORTED!
	Background BackgroundOptions
}

func (opts BracketOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Background.ToArgs(nil)...)

	return args
}
