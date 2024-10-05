package sketchybar

type ItemOptions struct {
	Icon        ItemIconOptions
	Label       ItemLabelOptions
	Background  BackgroundOptions
	Border      BorderOptions
	Padding     PaddingOptions
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

func (opts ItemOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Background.ToArgs(nil)...)
	args = append(args, opts.Label.ToArgs()...)
	args = append(args, opts.Icon.ToArgs()...)
	args = append(args, opts.Border.ToArgs(nil)...)
	args = append(args, opts.Padding.ToArgs(nil)...)

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

type ItemIconOptions struct {
	Padding    PaddingOptions
	Color      ColorOptions
	Background BackgroundOptions
	Font       FontOptions
	Drawing    string
	Value      string
	Highlight  string
}

func (opts ItemIconOptions) ToArgs() []string {
	args := []string{}

	parent := "icon"
	args = append(args, opts.Padding.ToArgs(&parent)...)
	args = append(args, opts.Color.ToArgs(&parent)...)
	args = append(args, opts.Background.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "icon=%s", opts.Value)
	}
	if opts.Font != EmptyFontOptions {
		args = with(args, "icon.font=%s", opts.Font.String())
	}
	if opts.Highlight != "" {
		args = with(args, "icon.highlight=%s", opts.Highlight)
	}
	if opts.Drawing != "" {
		args = with(args, "icon.drawing=%s", opts.Drawing)
	}

	return args
}

type ItemLabelOptions struct {
	Padding   PaddingOptions
	Color     ColorOptions
	Font      FontOptions
	Drawing   string
	Value     string
	Highlight string
}

func (opts ItemLabelOptions) ToArgs() []string {
	args := []string{}

	parent := "label"

	args = append(args, opts.Padding.ToArgs(&parent)...)
	args = append(args, opts.Color.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "label=%s", opts.Value)
	}
	if opts.Font != EmptyFontOptions {
		args = with(args, "label.font=%s", opts.Font.String())
	}
	if opts.Highlight != "" {
		args = with(args, "label.highlight=%s", opts.Highlight)
	}
	if opts.Drawing != "" {
		args = with(args, "label.drawing=%s", opts.Drawing)
	}
	return args
}
