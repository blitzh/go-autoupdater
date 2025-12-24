package updater

import (
	"strconv"
	"strings"
)

func CompareVersion(a, b string) int {
	// returns: -1 if a<b, 0 if equal, +1 if a>b
	pa := parseVersion(a)
	pb := parseVersion(b)

	n := len(pa)
	if len(pb) > n {
		n = len(pb)
	}

	for i := 0; i < n; i++ {
		ai := 0
		bi := 0
		if i < len(pa) {
			ai = pa[i]
		}
		if i < len(pb) {
			bi = pb[i]
		}
		if ai < bi {
			return -1
		}
		if ai > bi {
			return 1
		}
	}

	// if numeric equal but raw differs, do stable tie-break
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func parseVersion(v string) []int {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			// non-numeric => stop parsing to avoid wrong compare
			break
		}
		out = append(out, n)
	}
	if len(out) == 0 {
		return []int{0}
	}
	return out
}
