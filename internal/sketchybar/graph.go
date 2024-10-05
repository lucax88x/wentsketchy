package sketchybar

type GraphOptions struct {
	Icon        ItemIconOptions
	Label       ItemLabelOptions
	Background  BackgroundOptions
	Border      BorderOptions
	Padding     PaddingOptions
	Graph       ItemGraphOptions
	Display     string
	Space       string
	Width       *int
	YOffset     *int
	UpdateFreq  *int
	Updates     string
	ScrollTexts string
	Script      string
	ClickScript string
	MachHelper  string
}

func (opts GraphOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Background.ToArgs(nil)...)
	args = append(args, opts.Label.ToArgs()...)
	args = append(args, opts.Icon.ToArgs()...)
	args = append(args, opts.Border.ToArgs(nil)...)
	args = append(args, opts.Padding.ToArgs(nil)...)
	args = append(args, opts.Graph.ToArgs()...)

	if opts.Display != "" {
		args = with(args, "display=%s", opts.Display)
	}
	if opts.Space != "" {
		args = with(args, "space=%s", opts.Space)
	}
	if opts.Width != nil {
		args = with(args, "width=%d", *opts.Width)
	}
	if opts.YOffset != nil {
		args = with(args, "y_offset=%d", *opts.YOffset)
	}
	if opts.UpdateFreq != nil {
		args = with(args, "update_freq=%d", *opts.UpdateFreq)
	}
	if opts.Updates != "" {
		args = with(args, "updates=%s", opts.Updates)
	}
	if opts.ScrollTexts != "" {
		args = with(args, "scroll_texts=%s", opts.ScrollTexts)
	}
	if opts.Script != "" {
		args = with(args, "script=%s", opts.Script)
	}
	if opts.ClickScript != "" {
		args = with(args, "click_script=%s", opts.ClickScript)
	}
	if opts.MachHelper != "" {
		args = with(args, "mach_helper=%s", opts.MachHelper)
	}

	return args
}

type ItemGraphOptions struct {
	Color     string
	FillColor string
}

func (opts ItemGraphOptions) ToArgs() []string {
	args := []string{}

	if opts.Color != "" {
		args = with(args, "graph.color=%s", opts.Color)
	}
	if opts.FillColor != "" {
		args = with(args, "graph.fill_color=%s", opts.FillColor)
	}
	return args
}
