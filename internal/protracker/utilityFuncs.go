package protracker

import (
	"strconv"
)

func uint8FromHexString(s string) (uint8, error) {
	u, err := strconv.ParseUint(s, 16, 8)
	return uint8(u), err
}

// duplicate. Figure out how to share between packages (symlink won't do it)
func compareSlices[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
			return false
		}
	}
	return true
}
