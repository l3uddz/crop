package stringutils

import "strings"

func FromLeftUntil(from string, sub string) string {
	pos := strings.Index(from, sub)
	if pos < 0 {
		return from
	}

	return from[:pos]
}
