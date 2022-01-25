package stringutil

// Contains returns true if `sl` contains `s`.
func Contains(s string, sl []string) bool {
	for i := range sl {
		if sl[i] == s {
			return true
		}
	}

	return false
}

// Remove removes `s` from `sl`
func Remove(s string, sl []string) []string {
	var out []string
	for i := range sl {
		if sl[i] == s {
			continue
		}

		out = append(out, sl[i])
	}
	return out
}
