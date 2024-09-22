package sketchybar

type BackgroundOptions struct {
	Border       BorderOptions
	Color        ColorOptions
	Image        ImageOptions
	Padding      PaddingOptions
	Drawing      string
	Width        int
	Height       int
	CornerRadius int
}

func (opts BackgroundOptions) ToArgs(parent *string) []string {
	args := []string{}

	parentAndPrefix := mergeParentAndPrefix(parent, "background")

	args = append(args, opts.Color.ToArgs(parentAndPrefix)...)
	args = append(args, opts.Border.ToArgs(parentAndPrefix)...)
	args = append(args, opts.Image.ToArgs(parentAndPrefix)...)
	args = append(args, opts.Padding.ToArgs(parentAndPrefix)...)

	if opts.Drawing != "" {
		args = withParent(args, parent, "background.drawing=%s", opts.Drawing)
	}
	if opts.Width != 0 {
		args = withParent(args, parent, "background.width=%d", opts.Width)
	}
	if opts.Height != 0 {
		args = withParent(args, parent, "background.height=%d", opts.Height)
	}
	if opts.CornerRadius != 0 {
		args = withParent(args, parent, "background.corner_radius=%d", opts.CornerRadius)
	}

	return args
}
