package items

type Batches = [][]string

func batch(arr Batches, args []string) Batches {
	return append(arr, args)
}

func s(args ...string) []string {
	return args
}

func m(left []string, right []string) []string {
	return append(left, right...)
}

func pointer(i int) *int {
	return &i
}

func Flatten(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
