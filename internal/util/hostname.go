package util

import "strings"

// TODO: comment exported function and code sections
func HostnamesMatch(hostnameA, hostnameB string) bool {
	// the hostnames are in the form of "foo.bar.com"; split them
	// in a slice of substrings
	hostnameAParts := strings.Split(hostnameA, ".")
	hostnameBParts := strings.Split(hostnameB, ".")

	var a, b, a2, b2 int
	var wildcardA, wildcardB bool

	// iterate over the parts of both the hostnames
	for a, b = 0, 0; a < len(hostnameAParts) && b < len(hostnameBParts); a, b = a+1, b+1 {
		var matchFoundA, matchFoundB bool

		// if the current part of B is a wildcard, we need to find the first
		// A part that matches with the following B part
		if wildcardB {
			for b2 = b; b2 < len(hostnameBParts); b2 += 1 {
				if hostnameAParts[a] == hostnameBParts[b2] {
					matchFoundA = true
					break
				}
			}
		}

		// if the current part of A is a wildcard, we need to find the first
		// B part that matches with the following A part
		if wildcardA {
			for a2 = a; a2 < len(hostnameAParts); a2 += 1 {
				if hostnameAParts[a2] == hostnameBParts[b] {
					matchFoundB = true
					break
				}
			}
		}

		if wildcardB || wildcardA {
			// if wither one of the parts was a wildcard and no match was found,
			//the hostnames are incompatible
			if !matchFoundA && !matchFoundB {
				return false
			}
			// set the b index to the new shifted value
			if matchFoundA {
				b = b2
			}
			// set the a index to the new shifted value
			if matchFoundB {
				a = a2
			}
		}

		// check if at least on of the current parts are a wildcard; if so, continue
		wildcardB = hostnameAParts[a] == "*"
		wildcardA = hostnameBParts[b] == "*"
		if wildcardB || wildcardA {
			continue
		}
		// reset the wildcard  variables
		wildcardB, wildcardA = false, false

		// if the current a part is different from the b part, the hostnames are incompatible
		if hostnameAParts[a] != hostnameBParts[b] {
			return false
		}
	}
	if len(hostnameBParts)-b != len(hostnameAParts)-a {
		return false
	}

	return true
}
