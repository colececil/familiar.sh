package packagemanagers

import (
	"strconv"
	"strings"
	"unicode"
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
	version1Parts := strings.FieldsFunc(versionString1, isRuneNonAlphanumeric)
	version2Parts := strings.FieldsFunc(versionString2, isRuneNonAlphanumeric)

	for i := 0; i < len(version1Parts) || i < len(version2Parts); i++ {
		version1Part := "0"
		if i < len(version1Parts) {
			version1Part = version1Parts[i]
		}

		version2Part := "0"
		if i < len(version2Parts) {
			version2Part = version2Parts[i]
		}

		bothAreNumbers := true

		version1Number, err := strconv.Atoi(version1Part)
		if err != nil {
			bothAreNumbers = false
		}

		version2Number, err := strconv.Atoi(version2Part)
		if err != nil {
			bothAreNumbers = false
		}

		if bothAreNumbers {
			if version1Number < version2Number {
				return -1
			} else if version1Number > version2Number {
				return 1
			}
		} else {
			if version1Part < version2Part {
				return -1
			} else if version1Part > version2Part {
				return 1
			}
		}
	}

	return 0
}

// isRuneNonAlphanumeric returns whether the rune is an alphanumeric character.
func isRuneNonAlphanumeric(character rune) bool {
	return !unicode.IsNumber(character) && !unicode.IsLetter(character)
}
