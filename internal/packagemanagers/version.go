package packagemanagers

import (
	"strings"
)

// Version represents a version of a package or package manager.
type Version struct {
	VersionString string
}

// String returns the string representation of the version.
func (version *Version) String() string {
	return version.VersionString
}

// IsEqualTo returns whether the version is equal to the other version.
func (version *Version) IsEqualTo(otherVersion *Version) bool {
	return compareVersionStrings(version.VersionString, otherVersion.VersionString) == 0
}

// IsLessThan returns whether the version is less than the other version.
func (version *Version) IsLessThan(otherVersion *Version) bool {
	return compareVersionStrings(version.VersionString, otherVersion.VersionString) == -1
}

// IsGreaterThan returns whether the version is greater than the other version.
func (version *Version) IsGreaterThan(otherVersion *Version) bool {
	return compareVersionStrings(version.VersionString, otherVersion.VersionString) == 1
}

// compareVersionStrings compares two version strings. It returns -1 if versionString1 is less than versionString2, 0 if
// versionString1 is equal to versionString2, and 1 if versionString1 is greater than versionString2.
func compareVersionStrings(versionString1 string, versionString2 string) int {
	versionString1Parts := strings.Split(versionString1, ".")
	versionString2Parts := strings.Split(versionString2, ".")

	for i := 0; i < len(versionString1Parts) && i < len(versionString2Parts); i++ {
		if versionString1Parts[i] < versionString2Parts[i] {
			return -1
		} else if versionString1Parts[i] > versionString2Parts[i] {
			return 1
		}
	}

	if len(versionString1Parts) == len(versionString2Parts) {
		return 0
	}

	if len(versionString1Parts) < len(versionString2Parts) {
		return -1
	} else {
		return 1
	}
}
