package stringutils

import (
	"fmt"
	"strings"
)

func LeftJust(text string, filler string, size int) string {
	repeatSize := size - len(text)
	return fmt.Sprintf("%s%s", text, strings.Repeat(filler, repeatSize))
}
