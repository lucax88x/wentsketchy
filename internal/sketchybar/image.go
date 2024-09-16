package sketchybar

type ImageOptions struct {
	BorderOptions
	ColorOptions
	Value   string
	Drawing bool
	Scale   string
}

func (opts ImageOptions) ToArgs(parent *string) []string {
	args := []string{}

	parentAndPrefix := mergeParentAndPrefix(parent, "image")

	args = append(args, opts.ColorOptions.ToArgs(parentAndPrefix)...)
	args = append(args, opts.BorderOptions.ToArgs(parentAndPrefix)...)

	if opts.Value != "" {
		args = withParent(args, parent, "image=%s", opts.Value)
	}

	if opts.Drawing {
		args = withParent(args, parent, "image.drawing=%s", "on")
	} else {
		args = withParent(args, parent, "image.drawing=%s", "off")
	}

	if opts.Scale != "" {
		args = withParent(args, parent, "image.scale=%s", opts.Scale)
	}

	return args
}
