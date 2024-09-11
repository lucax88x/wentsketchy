package sketchybar

type SpaceOptions struct {
	ItemOptions
}

func (opts SpaceOptions) ToArgs() []string {
	args := []string{}

	args = append(args, opts.ItemOptions.ToArgs()...)

	return args
}
