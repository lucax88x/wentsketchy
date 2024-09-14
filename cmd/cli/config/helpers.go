package config

func flatten(slices ...[]string) []string {
	result := []string{}
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
