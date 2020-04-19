package rclone

import "strings"

func FormattedParams(params []string) []string {
	var formattedParams []string

	for _, p := range params {
		// does p contain =
		pos := strings.Index(p, "=")
		if pos < 0 {
			// no = found, add it as-is
			formattedParams = append(formattedParams, p)
			continue
		}

		// split p into param name and value
		pn := p[:pos]
		pv := p[pos+1:]

		if pn == "" || pv == "" {
			// failed to parse name and value
			log.Warnf("Failed formatting argument into name / value: %q, ignoring...", p)
			continue
		}

		formattedParams = append(formattedParams, pn, pv)
	}

	return formattedParams
}
