package util

import "strings"

// TODO: comment exported function and code sections
func Match(pattern, value string) bool {
	patternParts := strings.Split(pattern, ".")
	valueParts := strings.Split(value, ".")

	var i, j int
	var wildcardj, wildcardi bool
	for i, j = 0, 0; i < len(patternParts) && j < len(valueParts); i, j = i+1, j+1 {
		var matchFound bool

		if wildcardj {
			for ; i < len(patternParts); i += 1 {
				if patternParts[i] == valueParts[j] {
					matchFound = true
					break
				}
			}
			if !matchFound {
				return false
			}
		}

		if wildcardi {
			for ; j < len(valueParts); j += 1 {
				if patternParts[i] == valueParts[j] {
					matchFound = true
					break
				}
			}
			if !matchFound {
				return false
			}
		}

		if patternParts[i] == "*" {
			wildcardi = true
			continue
		}
		if valueParts[j] == "*" {
			wildcardj = true
			continue
		}

		wildcardi, wildcardj = false, false

		if patternParts[i] != valueParts[j] {
			return false
		}
	}
	if len(valueParts)-j != len(patternParts)-i {
		return false
	}

	return true
}
