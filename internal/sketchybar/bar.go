package sketchybar

type BarOptions struct {
	PaddingOptions
	ColorOptions
	BorderOptions
	Height     int
	Shadow     bool
	Background BackgroundOptions
	Position   string
	Sticky     bool
	YOffset    int
	Margin     int
	Topmost    string
}

func (opts BarOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.PaddingOptions.ToArgs(nil)...)
	args = append(args, opts.ColorOptions.ToArgs(nil)...)
	args = append(args, opts.BorderOptions.ToArgs(nil)...)

	if opts.Height != 0 {
		args = with(args, "height=%d", opts.Height)
	}
	if opts.Shadow {
		args = with(args, "shadow", "on")
	}
	if opts.Position != "" {
		args = with(args, "position=%s", opts.Position)
	}
	if opts.Sticky {
		args = with(args, "sticky=%s", "on")
	}
	if opts.YOffset != 0 {
		args = with(args, "y_offset=%d", opts.YOffset)
	}
	if opts.Margin != 0 {
		args = with(args, "margin=%d", opts.Margin)
	}
	if opts.Topmost != "" {
		args = with(args, "topmost=%s", opts.Topmost)
	}

	return args
}
