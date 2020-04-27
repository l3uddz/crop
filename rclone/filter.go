package rclone

import "fmt"

func IncludeExcludeToFilters(includes []string, excludes []string) []string {
	params := make([]string, 0)

	// add excludes
	if len(excludes) > 0 {
		for _, exclude := range excludes {
			params = append(params, "--filter", fmt.Sprintf("- %s", exclude))
		}
	}

	// were there includes?
	if len(includes) > 0 {
		for _, include := range includes {
			params = append(params, "--filter", fmt.Sprintf("+ %s", include))
		}

		// includes need the below, see: https://forum.rclone.org/t/filter-or-include-exclude-help-needed/10890/2
		params = append(params, "--filter", "- *")
	}

	return params
}
