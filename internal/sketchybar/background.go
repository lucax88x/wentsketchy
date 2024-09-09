package sketchybar

type BackgroundOptions struct {
	BorderOptions
	ColorOptions
}

func (opts BackgroundOptions) ToArgs() []string {
	args := []string{}

	prefix := "background"

	args = append(args, opts.ColorOptions.ToArgs(&prefix)...)
	args = append(args, opts.BorderOptions.ToArgs(&prefix)...)

	return args
}
