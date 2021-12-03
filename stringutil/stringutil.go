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
