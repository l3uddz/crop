package reutils

import (
	"regexp"
	"strings"
)

var (
	num = regexp.MustCompile(`(\d+)`)
)

func GetEveryNumber(from string) string {
	matches := num.FindAllString(from, -1)
	val := strings.Join(matches, "")
	return strings.TrimLeft(val, "0")
}
