package sketchybar

type BarOptions struct {
	Padding       PaddingOptions
	Color         ColorOptions
	Border        BorderOptions
	Height        *int
	Shadow        string
	FontSmoothing string
	Background    BackgroundOptions
	Position      string
	Sticky        string
	YOffset       *int
	Margin        *int
	Topmost       string
}

func (opts BarOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Padding.ToArgs(nil)...)
	args = append(args, opts.Color.ToArgs(nil)...)
	args = append(args, opts.Border.ToArgs(nil)...)

	if opts.Height != nil {
		args = with(args, "height=%d", *opts.Height)
	}
	if opts.Shadow != "" {
		args = with(args, "shadow=%s", opts.Shadow)
	}
	if opts.Position != "" {
		args = with(args, "position=%s", opts.Position)
	}
	if opts.Sticky != "" {
		args = with(args, "sticky=%s", opts.Sticky)
	}
	if opts.FontSmoothing != "" {
		args = with(args, "font_smoothing=%s", opts.FontSmoothing)
	}
	if opts.YOffset != nil {
		args = with(args, "y_offset=%d", *opts.YOffset)
	}
	if opts.Margin != nil {
		args = with(args, "margin=%d", *opts.Margin)
	}
	if opts.Topmost != "" {
		args = with(args, "topmost=%s", opts.Topmost)
	}

	return args
}
