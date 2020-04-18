package sliceutils

import "strings"

func StringSliceContains(slice []string, contains string, caseInsensitive bool) bool {
	for _, str := range slice {
		match := false

		switch caseInsensitive {
		case true:
			match = strings.EqualFold(str, contains)
		default:
			match = str == contains
		}

		if match {
			// slice contained the string
			return true
		}

	}

	return false
}
