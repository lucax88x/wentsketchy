package items

func batch(arr [][]string, args []string) [][]string {
	return append(arr, args)
}

func s(args ...string) []string {
	return args
}

func m(left []string, right []string) []string {
	return append(left, right...)
}
