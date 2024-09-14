package sketchybar

type BackgroundOptions struct {
	BorderOptions
	ColorOptions
	ImageOptions
	Drawing bool
}

func (opts BackgroundOptions) ToArgs(parent *string) []string {
	args := []string{}

	parentAndPrefix := mergeParentAndPrefix(parent, "background")

	args = append(args, opts.ColorOptions.ToArgs(parentAndPrefix)...)
	args = append(args, opts.BorderOptions.ToArgs(parentAndPrefix)...)
	args = append(args, opts.ImageOptions.ToArgs(parentAndPrefix)...)

	if opts.Drawing {
		args = withParent(args, parent, "background.drawing=%s", "on")
	}

	return args
}
