package packagemanagers_test

import (
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	Describe("String", func() {
		It("should return the `VersionString` field", func() {
			versionString := "versionString"
			version := NewVersion(versionString)

			result := version.String()
			Expect(result).To(Equal(versionString))
		})
	})

	Describe("IsEqualTo", func() {
		It("should return true when the version strings are equal", func() {
			versionString := "1.0.0"
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)

			result := version1.IsEqualTo(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when the version strings are not equal", func() {
			version1 := NewVersion("1.0.0")
			version2 := NewVersion("1.0.1")

			result := version1.IsEqualTo(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings are empty", func() {
			versionString := ""
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)
			result := version1.IsEqualTo(version2)
			Expect(result).To(BeTrue())
		})

		It("should handle sections with only numbers, sections with only non-numeric characters, and sections with a "+
			"mix of both", func() {
			versionString := "123.abc.123abc"
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)

			result := version1.IsEqualTo(version2)
			Expect(result).To(BeTrue())
		})

		It("should consider missing sections to be 0 when they are present in one string but not the other", func() {
			version1 := NewVersion("1.0")
			version2 := NewVersion("1.0.0")
			version3 := NewVersion("1.0.1")
			version4 := NewVersion("1.0.abc")

			result := version1.IsEqualTo(version2)
			Expect(result).To(BeTrue())

			result = version1.IsEqualTo(version3)
			Expect(result).To(BeFalse())

			result = version1.IsEqualTo(version4)
			Expect(result).To(BeFalse())
		})

		It("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
			"different separators", func() {
			version1 := NewVersion("1.2-3_4:5+6~7")
			version2 := NewVersion("1~2.3-4_5:6+7")

			result := version1.IsEqualTo(version2)
			Expect(result).To(BeTrue())
		})
	})

	Describe("IsLessThan", func() {
		It("should return true when both version strings have one part and the first one is less than the second",
			func() {
				version1 := NewVersion("1")
				version2 := NewVersion("2")

				result := version1.IsLessThan(version2)
				Expect(result).To(BeTrue())
			})

		It("should return false when both version strings have one part and the first one is greater than the second",
			func() {
				version1 := NewVersion("2")
				version2 := NewVersion("1")

				result := version1.IsLessThan(version2)
				Expect(result).To(BeFalse())
			})

		It("should return false when both version strings have one part and the first one is equal to the second",
			func() {
				versionString := "1"
				version1 := NewVersion(versionString)
				version2 := NewVersion(versionString)

				result := version1.IsLessThan(version2)
				Expect(result).To(BeFalse())
			})

		It("should return false when both version strings are empty", func() {
			versionString := ""
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)
			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings have multiple parts and the middle part of the first one is "+
			"less than that of the second", func() {
			version1 := NewVersion("1.0.0")
			version2 := NewVersion("1.1.0")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when both version strings have multiple parts and the middle part of the first one is "+
			"greater than that of the second", func() {
			version1 := NewVersion("1.1.0")
			version2 := NewVersion("1.0.0")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings have multiple parts and the final part of the first one is "+
			"less than that of the second", func() {
			version1 := NewVersion("1.0.0")
			version2 := NewVersion("1.0.1")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when both version strings have multiple parts and the final part of the first one is "+
			"greater than that of the second", func() {
			version1 := NewVersion("1.0.1")
			version2 := NewVersion("1.0.0")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings have multiple parts and the first part of the first string "+
			"is less than that of the second, even if the first string is not less in later parts", func() {
			version1 := NewVersion("1.1.1")
			version2 := NewVersion("2.0.0")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when both version strings have multiple parts and the first part of the first string "+
			"is greater than that of the second, even if the first string is less in later parts", func() {
			version1 := NewVersion("2.0.0")
			version2 := NewVersion("1.1.1")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return false when both version strings have multiple parts and all parts are equal", func() {
			versionString := "1.0.0"
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should compare version string sections numerically when they only contain numbers", func() {
			version1 := NewVersion("9")
			version2 := NewVersion("10")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should compare version string sections lexicographically when they contain only non-numeric characters",
			func() {
				version1 := NewVersion("abc")
				version2 := NewVersion("acc")
				version3 := NewVersion("abcd")

				result := version1.IsLessThan(version2)
				Expect(result).To(BeTrue())

				result = version1.IsLessThan(version3)
				Expect(result).To(BeTrue())
			})

		It("should compare version string sections lexicographically when they contain both numbers and other "+
			"characters", func() {
			version1 := NewVersion("9a")
			version2 := NewVersion("10a")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should consider missing sections to be 0 when they are present in one string but not the other", func() {
			version1 := NewVersion("1.0")
			version2 := NewVersion("1.0.0")
			version3 := NewVersion("1.0.1")
			version4 := NewVersion("1.0.abc")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())

			result = version2.IsLessThan(version1)
			Expect(result).To(BeFalse())

			result = version1.IsLessThan(version3)
			Expect(result).To(BeTrue())

			result = version3.IsLessThan(version1)
			Expect(result).To(BeFalse())

			result = version1.IsLessThan(version4)
			Expect(result).To(BeTrue())

			result = version4.IsLessThan(version1)
			Expect(result).To(BeFalse())
		})

		It("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
			"different separators", func() {
			version1 := NewVersion("1.2-3_4:5+6~7")
			version2 := NewVersion("1~2.3-4_5:6+6")
			version3 := NewVersion("1~2.3-4_5:6+8")

			result := version1.IsLessThan(version2)
			Expect(result).To(BeFalse())

			result = version1.IsLessThan(version3)
			Expect(result).To(BeTrue())
		})
	})

	Describe("IsGreaterThan", func() {
		It("should return false when both version strings have one part and the first one is less than the second",
			func() {
				version1 := NewVersion("1")
				version2 := NewVersion("2")

				result := version1.IsGreaterThan(version2)
				Expect(result).To(BeFalse())
			})

		It("should return true when both version strings have one part and the first one is greater than the second",
			func() {
				version1 := NewVersion("2")
				version2 := NewVersion("1")

				result := version1.IsGreaterThan(version2)
				Expect(result).To(BeTrue())
			})

		It("should return false when both version strings have one part and the first one is equal to the second",
			func() {
				versionString := "1"
				version1 := NewVersion(versionString)
				version2 := NewVersion(versionString)

				result := version1.IsGreaterThan(version2)
				Expect(result).To(BeFalse())
			})

		It("should return false when both version strings are empty", func() {
			versionString := ""
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)
			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return false when both version strings have multiple parts and the middle part of the first one is "+
			"less than that of the second", func() {
			version1 := NewVersion("1.0.0")
			version2 := NewVersion("1.1.0")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings have multiple parts and the middle part of the first one is "+
			"greater than that of the second", func() {
			version1 := NewVersion("1.1.0")
			version2 := NewVersion("1.0.0")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when both version strings have multiple parts and the final part of the first one is "+
			"less than that of the second", func() {
			version1 := NewVersion("1.0.0")
			version2 := NewVersion("1.0.1")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return true when both version strings have multiple parts and the final part of the first one is "+
			"greater than that of the second", func() {
			version1 := NewVersion("1.0.1")
			version2 := NewVersion("1.0.0")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return true when both version strings have multiple parts and the first part of the first string "+
			"is greater than that of the second, even if the first string is not greater in later parts", func() {
			version1 := NewVersion("2.0.0")
			version2 := NewVersion("1.1.1")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should return false when both version strings have multiple parts and the first part of the first string "+
			"is less than that of the second, even if the first string is greater in later parts", func() {
			version1 := NewVersion("1.1.1")
			version2 := NewVersion("2.0.0")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should return false when both version strings have multiple parts and all parts are equal", func() {
			versionString := "1.0.0"
			version1 := NewVersion(versionString)
			version2 := NewVersion(versionString)

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())
		})

		It("should compare version string sections numerically when they only contain numbers", func() {
			version1 := NewVersion("10")
			version2 := NewVersion("9")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should compare version string sections lexicographically when they contain only non-numeric characters",
			func() {
				version1 := NewVersion("abc")
				version2 := NewVersion("acc")
				version3 := NewVersion("abcd")

				result := version1.IsGreaterThan(version2)
				Expect(result).To(BeFalse())

				result = version1.IsGreaterThan(version3)
				Expect(result).To(BeFalse())
			})

		It("should compare version string sections lexicographically when they contain both numbers and other "+
			"characters", func() {
			version1 := NewVersion("9a")
			version2 := NewVersion("10a")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())
		})

		It("should consider missing sections to be 0 when they are present in one string but not the other", func() {
			version1 := NewVersion("1.0")
			version2 := NewVersion("1.0.0")
			version3 := NewVersion("1.0.1")
			version4 := NewVersion("1.0.abc")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeFalse())

			result = version2.IsGreaterThan(version1)
			Expect(result).To(BeFalse())

			result = version1.IsGreaterThan(version3)
			Expect(result).To(BeFalse())

			result = version3.IsGreaterThan(version1)
			Expect(result).To(BeTrue())

			result = version1.IsGreaterThan(version4)
			Expect(result).To(BeFalse())

			result = version4.IsGreaterThan(version1)
			Expect(result).To(BeTrue())
		})

		It("should consider any non-alphanumeric character to be a separator, and it should not differentiate between "+
			"different separators", func() {
			version1 := NewVersion("1.2-3_4:5+6~7")
			version2 := NewVersion("1~2.3-4_5:6+6")
			version3 := NewVersion("1~2.3-4_5:6+8")

			result := version1.IsGreaterThan(version2)
			Expect(result).To(BeTrue())

			result = version1.IsGreaterThan(version3)
			Expect(result).To(BeFalse())
		})
	})
})
