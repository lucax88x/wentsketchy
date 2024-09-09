package config

func s(args ...string) []string {
	return args
}

func m(left []string, right []string) []string {
	return append(left, right...)
}
