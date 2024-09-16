package sketchybar

type BackgroundOptions struct {
	BorderOptions
	ColorOptions
	ImageOptions
	Drawing      bool
	Width        int
	Height       int
	CornerRadius int
}

func (opts BackgroundOptions) ToArgs(parent *string) []string {
	args := []string{}

	parentAndPrefix := mergeParentAndPrefix(parent, "background")

	args = append(args, opts.ColorOptions.ToArgs(parentAndPrefix)...)
	args = append(args, opts.BorderOptions.ToArgs(parentAndPrefix)...)
	args = append(args, opts.ImageOptions.ToArgs(parentAndPrefix)...)

	if opts.Drawing {
		args = withParent(args, parent, "background.drawing=%s", "on")
	} else {
		args = withParent(args, parent, "background.drawing=%s", "off")
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
