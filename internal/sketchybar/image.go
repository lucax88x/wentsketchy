package sketchybar

type ImageOptions struct {
	Border  BorderOptions
	Color   ColorOptions
	Padding PaddingOptions

	Value   string
	Drawing string
	Scale   string
}

func (opts ImageOptions) ToArgs(parent *string) []string {
	args := []string{}

	parentAndPrefix := mergeParentAndPrefix(parent, "image")

	args = append(args, opts.Border.ToArgs(parentAndPrefix)...)
	args = append(args, opts.Color.ToArgs(parentAndPrefix)...)
	args = append(args, opts.Padding.ToArgs(parentAndPrefix)...)

	if opts.Value != "" {
		args = withParent(args, parent, "image=%s", opts.Value)
	}

	if opts.Drawing != "" {
		args = withParent(args, parent, "image.drawing=%s", opts.Drawing)
	}

	if opts.Scale != "" {
		args = withParent(args, parent, "image.scale=%s", opts.Scale)
	}

	return args
}
