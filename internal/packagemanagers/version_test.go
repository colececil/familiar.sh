package packagemanagers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	t.Run("String", testString)
	t.Run("IsEqualTo", testIsEqualTo)
	t.Run("IsLessThan", testIsLessThan)
	t.Run("IsGreaterThan", testIsGreaterThan)
}

func testString(t *testing.T) {
	t.Run("should return the `VersionString` field", func(t *testing.T) {
		versionString := "versionString"
		version := Version{VersionString: versionString}

		result := version.String()
		assert.Equal(t, result, versionString)
	})
}

func testIsEqualTo(t *testing.T) {
	t.Run("should return true when the version strings are equal", func(t *testing.T) {
		versionString := "1.0.0"
		version1 := &Version{VersionString: versionString}
		version2 := &Version{VersionString: versionString}

		result := version1.IsEqualTo(version2)
		assert.True(t, result)
	})

	t.Run("should return false when the version strings are not equal", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.0"}
		version2 := &Version{VersionString: "1.0.1"}

		result := version1.IsEqualTo(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings are empty", func(t *testing.T) {
		versionString := ""
		version1 := &Version{VersionString: versionString}
		version2 := &Version{VersionString: versionString}
		result := version1.IsEqualTo(version2)
		assert.True(t, result)
	})

	t.Run("should handle sections with only numbers, sections with only non-numeric characters, and sections with a "+
		"mix of both", func(t *testing.T) {
		versionString := "123.abc.123abc"
		version1 := &Version{VersionString: versionString}
		version2 := &Version{VersionString: versionString}

		result := version1.IsEqualTo(version2)
		assert.True(t, result)
	})

	t.Run("should consider missing sections to be 0 when they are present in one string but not the other",
		func(t *testing.T) {
			version1 := &Version{VersionString: "1.0"}
			version2 := &Version{VersionString: "1.0.0"}
			version3 := &Version{VersionString: "1.0.1"}
			version4 := &Version{VersionString: "1.0.abc"}

			result := version1.IsEqualTo(version2)
			assert.True(t, result)

			result = version1.IsEqualTo(version3)
			assert.False(t, result)

			result = version1.IsEqualTo(version4)
			assert.False(t, result)
		})

	t.Run("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
		"different separators", func(t *testing.T) {
		version1 := &Version{VersionString: "1.2-3_4:5+6~7"}
		version2 := &Version{VersionString: "1~2.3-4_5:6+7"}

		result := version1.IsEqualTo(version2)
		assert.True(t, result)
	})
}

func testIsLessThan(t *testing.T) {
	t.Run("should return true when both version strings have one part and the first one is less than the second",
		func(t *testing.T) {
			version1 := &Version{VersionString: "1"}
			version2 := &Version{VersionString: "2"}

			result := version1.IsLessThan(version2)
			assert.True(t, result)
		})

	t.Run("should return false when both version strings have one part and the first one is greater than the second",
		func(t *testing.T) {
			version1 := &Version{VersionString: "2"}
			version2 := &Version{VersionString: "1"}

			result := version1.IsLessThan(version2)
			assert.False(t, result)
		})

	t.Run("should return false when both version strings have one part and the first one is equal to the second",
		func(t *testing.T) {
			versionString := "1"
			version1 := &Version{VersionString: versionString}
			version2 := &Version{VersionString: versionString}

			result := version1.IsLessThan(version2)
			assert.False(t, result)
		})

	t.Run("should return false when both version strings are empty", func(t *testing.T) {
		versionString := ""
		version1 := &Version{VersionString: versionString}
		version2 := &Version{VersionString: versionString}
		result := version1.IsLessThan(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the middle part of the first one is "+
		"less than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.0"}
		version2 := &Version{VersionString: "1.1.0"}

		result := version1.IsLessThan(version2)
		assert.True(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the middle part of the first one is "+
		"greater than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.1.0"}
		version2 := &Version{VersionString: "1.0.0"}

		result := version1.IsLessThan(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the final part of the first one is "+
		"less than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.0"}
		version2 := &Version{VersionString: "1.0.1"}

		result := version1.IsLessThan(version2)
		assert.True(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the final part of the first one is "+
		"greater than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.1"}
		version2 := &Version{VersionString: "1.0.0"}

		result := version1.IsLessThan(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the first part of the first string is "+
		"less than that of the second, even if the first string is not less in later parts", func(t *testing.T) {
		version1 := &Version{VersionString: "1.1.1"}
		version2 := &Version{VersionString: "2.0.0"}

		result := version1.IsLessThan(version2)
		assert.True(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the first part of the first string is "+
		"greater than that of the second, even if the first string is less in later parts", func(t *testing.T) {
		version1 := &Version{VersionString: "2.0.0"}
		version2 := &Version{VersionString: "1.1.1"}

		result := version1.IsLessThan(version2)
		assert.False(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and all parts are equal",
		func(t *testing.T) {
			versionString := "1.0.0"
			version1 := &Version{VersionString: versionString}
			version2 := &Version{VersionString: versionString}

			result := version1.IsLessThan(version2)
			assert.False(t, result)
		})

	t.Run("should compare version string sections numerically when they only contain numbers", func(t *testing.T) {
		version1 := &Version{VersionString: "9"}
		version2 := &Version{VersionString: "10"}

		result := version1.IsLessThan(version2)
		assert.True(t, result)
	})

	t.Run("should compare version string sections lexicographically when they contain only non-numeric characters",
		func(t *testing.T) {
			version1 := &Version{VersionString: "abc"}
			version2 := &Version{VersionString: "acc"}
			version3 := &Version{VersionString: "abcd"}

			result := version1.IsLessThan(version2)
			assert.True(t, result)

			result = version1.IsLessThan(version3)
			assert.True(t, result)
		})

	t.Run("should compare version string sections lexicographically when they contain both numbers and other "+
		"characters", func(t *testing.T) {
		version1 := &Version{VersionString: "9a"}
		version2 := &Version{VersionString: "10a"}

		result := version1.IsLessThan(version2)
		assert.False(t, result)
	})

	t.Run("should consider missing sections to be 0 when they are present in one string but not the other",
		func(t *testing.T) {
			version1 := &Version{VersionString: "1.0"}
			version2 := &Version{VersionString: "1.0.0"}
			version3 := &Version{VersionString: "1.0.1"}
			version4 := &Version{VersionString: "1.0.abc"}

			result := version1.IsLessThan(version2)
			assert.False(t, result)

			result = version2.IsLessThan(version1)
			assert.False(t, result)

			result = version1.IsLessThan(version3)
			assert.True(t, result)

			result = version3.IsLessThan(version1)
			assert.False(t, result)

			result = version1.IsLessThan(version4)
			assert.True(t, result)

			result = version4.IsLessThan(version1)
			assert.False(t, result)
		})

	t.Run("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
		"different separators", func(t *testing.T) {
		version1 := &Version{VersionString: "1.2-3_4:5+6~7"}
		version2 := &Version{VersionString: "1~2.3-4_5:6+6"}
		version3 := &Version{VersionString: "1~2.3-4_5:6+8"}

		result := version1.IsLessThan(version2)
		assert.False(t, result)

		result = version1.IsLessThan(version3)
		assert.True(t, result)
	})
}

func testIsGreaterThan(t *testing.T) {
	t.Run("should return false when both version strings have one part and the first one is less than the second",
		func(t *testing.T) {
			version1 := &Version{VersionString: "1"}
			version2 := &Version{VersionString: "2"}

			result := version1.IsGreaterThan(version2)
			assert.False(t, result)
		})

	t.Run("should return true when both version strings have one part and the first one is greater than the second",
		func(t *testing.T) {
			version1 := &Version{VersionString: "2"}
			version2 := &Version{VersionString: "1"}

			result := version1.IsGreaterThan(version2)
			assert.True(t, result)
		})

	t.Run("should return false when both version strings have one part and the first one is equal to the second",
		func(t *testing.T) {
			versionString := "1"
			version1 := &Version{VersionString: versionString}
			version2 := &Version{VersionString: versionString}

			result := version1.IsGreaterThan(version2)
			assert.False(t, result)
		})

	t.Run("should return false when both version strings are empty", func(t *testing.T) {
		versionString := ""
		version1 := &Version{VersionString: versionString}
		version2 := &Version{VersionString: versionString}
		result := version1.IsGreaterThan(version2)
		assert.False(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the middle part of the first one is "+
		"less than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.0"}
		version2 := &Version{VersionString: "1.1.0"}

		result := version1.IsGreaterThan(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the middle part of the first one is "+
		"greater than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.1.0"}
		version2 := &Version{VersionString: "1.0.0"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the final part of the first one is "+
		"less than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.0"}
		version2 := &Version{VersionString: "1.0.1"}

		result := version1.IsGreaterThan(version2)
		assert.False(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the final part of the first one is "+
		"greater than that of the second", func(t *testing.T) {
		version1 := &Version{VersionString: "1.0.1"}
		version2 := &Version{VersionString: "1.0.0"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)
	})

	t.Run("should return true when both version strings have multiple parts and the first part of the first string is "+
		"greater than that of the second, even if the first string is not greater in later parts", func(t *testing.T) {
		version1 := &Version{VersionString: "2.0.0"}
		version2 := &Version{VersionString: "1.1.1"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and the first part of the first string is "+
		"less than that of the second, even if the first string is greater in later parts", func(t *testing.T) {
		version1 := &Version{VersionString: "1.1.1"}
		version2 := &Version{VersionString: "2.0.0"}

		result := version1.IsGreaterThan(version2)
		assert.False(t, result)
	})

	t.Run("should return false when both version strings have multiple parts and all parts are equal",
		func(t *testing.T) {
			versionString := "1.0.0"
			version1 := &Version{VersionString: versionString}
			version2 := &Version{VersionString: versionString}

			result := version1.IsGreaterThan(version2)
			assert.False(t, result)
		})

	t.Run("should compare version string sections numerically when they only contain numbers", func(t *testing.T) {
		version1 := &Version{VersionString: "10"}
		version2 := &Version{VersionString: "9"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)
	})

	t.Run("should compare version string sections lexicographically when they contain only non-numeric characters",
		func(t *testing.T) {
			version1 := &Version{VersionString: "abc"}
			version2 := &Version{VersionString: "acc"}
			version3 := &Version{VersionString: "abcd"}

			result := version1.IsGreaterThan(version2)
			assert.False(t, result)

			result = version1.IsGreaterThan(version3)
			assert.False(t, result)
		})

	t.Run("should compare version string sections lexicographically when they contain both numbers and other "+
		"characters", func(t *testing.T) {
		version1 := &Version{VersionString: "9a"}
		version2 := &Version{VersionString: "10a"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)
	})

	t.Run("should consider missing sections to be 0 when they are present in one string but not the other",
		func(t *testing.T) {
			version1 := &Version{VersionString: "1.0"}
			version2 := &Version{VersionString: "1.0.0"}
			version3 := &Version{VersionString: "1.0.1"}
			version4 := &Version{VersionString: "1.0.abc"}

			result := version1.IsGreaterThan(version2)
			assert.False(t, result)

			result = version2.IsGreaterThan(version1)
			assert.False(t, result)

			result = version1.IsGreaterThan(version3)
			assert.False(t, result)

			result = version3.IsGreaterThan(version1)
			assert.True(t, result)

			result = version1.IsGreaterThan(version4)
			assert.False(t, result)

			result = version4.IsGreaterThan(version1)
			assert.True(t, result)
		})

	t.Run("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
		"different separators", func(t *testing.T) {
		version1 := &Version{VersionString: "1.2-3_4:5+6~7"}
		version2 := &Version{VersionString: "1~2.3-4_5:6+6"}
		version3 := &Version{VersionString: "1~2.3-4_5:6+8"}

		result := version1.IsGreaterThan(version2)
		assert.True(t, result)

		result = version1.IsGreaterThan(version3)
		assert.False(t, result)
	})
}
