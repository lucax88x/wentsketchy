package utils

import "strings"

func Sanitize(str string) string {
	return strings.TrimSpace(str)
}
