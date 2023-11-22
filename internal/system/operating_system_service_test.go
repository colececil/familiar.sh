package system_test

import (
	. "github.com/colececil/familiar.sh/internal/system"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OperatingSystemService", func() {
	var operatingSystemService *OperatingSystemService

	Describe("IsWindows", func() {
		It("should return true when the current operating system is Windows", func() {
			operatingSystemService = NewOperatingSystemService(WindowsOperatingSystem)
			result := operatingSystemService.IsWindows()
			Expect(result).To(BeTrue())
		})

		It("should return false when the current operating system is macOS", func() {
			operatingSystemService = NewOperatingSystemService(MacOsOperatingSystem)
			result := operatingSystemService.IsWindows()
			Expect(result).To(BeFalse())
		})

		It("should return false when the current operating system is Linux", func() {
			operatingSystemService = NewOperatingSystemService(LinuxOperatingSystem)
			result := operatingSystemService.IsWindows()
			Expect(result).To(BeFalse())
		})
	})
})
