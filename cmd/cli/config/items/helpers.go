package items

type batches = [][]string

func batch(arr batches, args []string) batches {
	return append(arr, args)
}

func s(args ...string) []string {
	return args
}

func m(left []string, right []string) []string {
	return append(left, right...)
}
