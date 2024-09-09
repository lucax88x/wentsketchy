package sketchybar

type ItemIconOptions struct {
	PaddingOptions
	Value string
	Font  FontOptions
}

func (opts ItemIconOptions) ToArgs() []string {
	args := []string{}

	parent := "icon"
	args = append(args, opts.PaddingOptions.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "icon=%s", opts.Value)
	}

	if opts.Font != EmptyFontOptions {
		args = with(args, "icon.font=%s'", opts.Font.String())
	}

	return args
}

type ItemLabelOptions struct {
	PaddingOptions
	Value string
}

func (opts ItemLabelOptions) ToArgs() []string {
	args := []string{}

	parent := "label"

	args = append(args, opts.PaddingOptions.ToArgs(&parent)...)

	if opts.Value != "" {
		args = with(args, "label=%s", opts.Value)
	}

	return args
}

type ItemOptions struct {
	Icon       ItemIconOptions
	Label      ItemLabelOptions
	Background BackgroundOptions
	Border     BorderOptions
	// BorderColor     string
	UpdateFreq  int
	Script      string
	ClickScript string
}

func (opts ItemOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.Background.ToArgs()...)
	args = append(args, opts.Label.ToArgs()...)
	args = append(args, opts.Icon.ToArgs()...)
	args = append(args, opts.Border.ToArgs(nil)...)

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
