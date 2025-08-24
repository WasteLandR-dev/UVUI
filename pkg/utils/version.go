package utils

import (
	"strconv"
	"strings"
)

// CompareVersions compares two version strings
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if v1 == v2
func CompareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int

		if i < len(parts1) {
			p1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			p2, _ = strconv.Atoi(parts2[i])
		}

		if p1 > p2 {
			return 1
		} else if p1 < p2 {
			return -1
		}
	}

	return 0
}

// IsNewerVersion checks if v1 is newer than v2
func IsNewerVersion(v1, v2 string) bool {
	return CompareVersions(v1, v2) > 0
}

// IsSameVersion checks if two versions are the same
func IsSameVersion(v1, v2 string) bool {
	return CompareVersions(v1, v2) == 0
}
