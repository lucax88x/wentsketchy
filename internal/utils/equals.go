package utils

import "strings"

func Equals(left string, right string) bool {
	return Sanitize(strings.ToLower(left)) == Sanitize(strings.ToLower(right))
}
