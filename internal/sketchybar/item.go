package sketchybar

type ItemIconOptions struct {
	PaddingOptions
	ColorOptions
	BackgroundOptions
	Value     string
	Font      FontOptions
	Highlight bool
}

func (opts ItemIconOptions) ToArgs() []string {
	args := []string{}

	parent := "icon"
	args = append(args, opts.PaddingOptions.ToArgs(&parent)...)
	args = append(args, opts.ColorOptions.ToArgs(&parent)...)
	args = append(args, opts.BackgroundOptions.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "icon=%s", opts.Value)
	}

	if opts.Font != EmptyFontOptions {
		args = with(args, "icon.font=%s", opts.Font.String())
	}

	if opts.Highlight {
		args = with(args, "icon.highlight=%s", "on")
	}

	return args
}

type ItemLabelOptions struct {
	PaddingOptions
	ColorOptions
	Value     string
	Font      FontOptions
	Highlight bool
}

func (opts ItemLabelOptions) ToArgs() []string {
	args := []string{}

	parent := "label"

	args = append(args, opts.PaddingOptions.ToArgs(&parent)...)
	args = append(args, opts.ColorOptions.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "label=%s", opts.Value)
	}
	if opts.Font != EmptyFontOptions {
		args = with(args, "label.font=%s", opts.Font.String())
	}
	if opts.Highlight {
		args = with(args, "icon.highlight=%s", "on")
	}

	return args
}

type ItemOptions struct {
	Icon        ItemIconOptions
	Label       ItemLabelOptions
	Background  BackgroundOptions
	Border      BorderOptions
	Display     string
	Space       string
	UpdateFreq  int
	Script      string
	ClickScript string
}

func (opts ItemOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Background.ToArgs(nil)...)
	args = append(args, opts.Label.ToArgs()...)
	args = append(args, opts.Icon.ToArgs()...)
	args = append(args, opts.Border.ToArgs(nil)...)

	if opts.Display != "" {
		args = with(args, "display=%s", opts.Display)
	}
	if opts.Space != "" {
		args = with(args, "space=%s", opts.Space)
	}
	if opts.UpdateFreq != 0 {
		args = with(args, "update_freq=%d", opts.UpdateFreq)
	}
	if opts.Script != "" {
		args = with(args, "script=%s", opts.Script)
	}
	if opts.ClickScript != "" {
		args = with(args, "click_script=%s", opts.ClickScript)
	}

	return args
}
