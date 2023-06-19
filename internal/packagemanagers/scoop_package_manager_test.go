package packagemanagers_test

import (
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/colececil/familiar.sh/internal/test"
)

var _ = Describe("ScoopPackageManager", func() {
	var operatingSystemServiceDouble *test.OperatingSystemServiceDouble
	var shellCommandServiceDouble *test.ShellCommandServiceDouble
	var scoopPackageManager *ScoopPackageManager

	BeforeEach(func() {
		operatingSystemServiceDouble = test.NewOperatingSystemServiceDouble()
		shellCommandServiceDouble = test.NewShellCommandServiceDouble()
		scoopPackageManager = NewScoopPackageManager(operatingSystemServiceDouble.OperatingSystemService,
			shellCommandServiceDouble.ShellCommandService)
	})

	Describe("Name", func() {
		It("should return \"scoop\"", func() {
			result := scoopPackageManager.Name()
			Expect(result).To(Equal("scoop"))
		})
	})

	Describe("IsSupported", func() {
	})

	Describe("IsInstalled", func() {
	})

	Describe("Install", func() {
	})

	Describe("Update", func() {
	})

	Describe("Uninstall", func() {
	})

	Describe("InstalledPackages", func() {
		var scoopExportOutput string
		var scoopStatusOutput string

		BeforeEach(func() {
			scoopExportOutput = `{
	"buckets": [
		{
			"Name": "main"
		}
	],
	"apps": [
		{
			"Source": "main",
			"Name": "package1",
			"Version": "1.0.0"
		},
		{
			"Source": "main",
			"Name": "package2",
			"Version": "2.3.4"
		},
		{
			"Source": "main",
			"Name": "package3",
			"Version": "3.2.1"
		}
	]
}
`
			scoopStatusOutput = `Scoop is up to date.

Name         Installed Version Latest Version   Missing Dependencies Info
----         ----------------- --------------   -------------------- ----
package1     1.0.0             1.0.0
package2     2.3.4             2.5.0
package3     3.2.1             4.0.0
`
		})

		It("should use the output of 'scoop export' to get the list of installed packages, along with the output of "+
			"'scoop status' to find out if there are newer package versions available", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopExportOutput, "scoop", false, "export")
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopStatusOutput, "scoop", false, "status")

			expectedPackages := []*Package{
				{
					Name:             "package1",
					InstalledVersion: &Version{VersionString: "1.0.0"},
					LatestVersion:    &Version{VersionString: "1.0.0"},
				},
				{
					Name:             "package2",
					InstalledVersion: &Version{VersionString: "2.3.4"},
					LatestVersion:    &Version{VersionString: "2.5.0"},
				},
				{
					Name:             "package3",
					InstalledVersion: &Version{VersionString: "3.2.1"},
					LatestVersion:    &Version{VersionString: "4.0.0"},
				},
			}

			packages, err := scoopPackageManager.InstalledPackages()
			Expect(err).To(BeNil())
			Expect(packages).To(Equal(expectedPackages))
		})

		PIt("should sort the packages by name", func() {
		})

		It("should return the correct information when all packages are up to date", func() {
			scoopStatusOutput := `Scoop is up to date.
Everything is ok!
`

			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopExportOutput, "scoop", false, "export")
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopStatusOutput, "scoop", false, "status")

			expectedPackages := []*Package{
				{
					Name:             "package1",
					InstalledVersion: &Version{VersionString: "1.0.0"},
					LatestVersion:    &Version{VersionString: "1.0.0"},
				},
				{
					Name:             "package2",
					InstalledVersion: &Version{VersionString: "2.3.4"},
					LatestVersion:    &Version{VersionString: "2.3.4"},
				},
				{
					Name:             "package3",
					InstalledVersion: &Version{VersionString: "3.2.1"},
					LatestVersion:    &Version{VersionString: "3.2.1"},
				},
			}

			packages, err := scoopPackageManager.InstalledPackages()
			Expect(err).To(BeNil())
			Expect(packages).To(Equal(expectedPackages))
		})
	})

	Describe("InstallPackage", func() {
	})

	Describe("UpdatePackage", func() {
	})

	Describe("UninstallPackage", func() {
	})
})
