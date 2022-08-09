package util

import "strings"

// TODO: comment exported function.
func HostnamesIntersect(hostnameA, hostnameB string) bool {
	return HostnamesMatch(hostnameA, hostnameB) || HostnamesMatch(hostnameB, hostnameA)
}

// TODO: comment exported function and code sections.
func HostnamesMatch(listenerHostname, routeHostname string) bool {
	// the hostnames are in the form of "foo.bar.com"; split them
	// in a slice of substrings
	listenerHostnameLabels := strings.Split(listenerHostname, ".")
	routeHostnameLabels := strings.Split(routeHostname, ".")

	var a, b int
	var wildcard bool

	// iterate over the parts of both the hostnames
	for a, b = 0, 0; a < len(listenerHostnameLabels) && b < len(routeHostnameLabels); a, b = a+1, b+1 {
		var matchFound bool

		// if the current part of B is a wildcard, we need to find the first
		// A part that matches with the following B part
		if wildcard {
			for ; b < len(routeHostnameLabels); b++ {
				if listenerHostnameLabels[a] == routeHostnameLabels[b] {
					matchFound = true
					break
				}
			}
		}

		if wildcard && !matchFound {
			return false
		}

		// check if at least on of the current parts are a wildcard; if so, continue
		if listenerHostnameLabels[a] == "*" {
			wildcard = true
			continue
		}
		// reset the wildcard  variables
		wildcard = false

		// if the current a part is different from the b part, the hostnames are incompatible
		if listenerHostnameLabels[a] != routeHostnameLabels[b] {
			return false
		}
	}
	return len(routeHostnameLabels)-b != len(listenerHostnameLabels)-a
}
