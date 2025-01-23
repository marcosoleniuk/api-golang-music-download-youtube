package utils

import "strings"

func IsValidFormat(format string) bool {
	return strings.Contains(SupportedFormats, format)
}
