package maputils

import (
	"fmt"
	"strings"
)

func GetStringMapValue(stringMap map[string]string, key string, caseSensitive bool) (string, error) {
	lowerKey := strings.ToLower(key)

	// case sensitive match
	if caseSensitive {
		v, ok := stringMap[key]
		if !ok {
			return "", fmt.Errorf("key was not found in map: %q", key)
		}

		return v, nil
	}

	// case insensitive match
	for k, v := range stringMap {
		if strings.ToLower(k) == lowerKey {
			return v, nil
		}
	}

	return "", fmt.Errorf("key was not found in map: %q", lowerKey)
}
