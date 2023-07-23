package packagemanagers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"

	"github.com/colececil/familiar.sh/internal/test"
)

var _ = Describe("ScoopPackageManager", func() {
	var operatingSystemServiceDouble *test.OperatingSystemServiceDouble
	var shellCommandServiceDouble *test.ShellCommandServiceDouble
	var outputWriterDouble *bytes.Buffer
	var scoopPackageManager *ScoopPackageManager

	BeforeEach(func() {
		operatingSystemServiceDouble = test.NewOperatingSystemServiceDouble()
		shellCommandServiceDouble = test.NewShellCommandServiceDouble()
		outputWriterDouble = new(bytes.Buffer)
		scoopPackageManager = NewScoopPackageManager(operatingSystemServiceDouble.OperatingSystemService,
			shellCommandServiceDouble.ShellCommandService, outputWriterDouble)
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

		JustBeforeEach(func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopExportOutput, "scoop", false, "export")
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopStatusOutput, "scoop", false, "status")
		})

		It("should write output stating it is getting installed package information", func() {
			scoopPackageManager.InstalledPackages()
			Expect(outputWriterDouble.String()).To(Equal(
				"Getting installed package information from package manager \"scoop\"...\n"))
		})

		When("all packages returned by `scoop export` are included in `scoop status`", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = createScoopStatusOutput([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				})
			})

			It("should indicate that all packages have updates available", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				}))
			})
		})

		When("the outputs of `scoop export` and `scoop status` are not ordered by package name", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package3", Version: "3.2.1", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package1", Version: "1.0.0", Source: "main"},
				})
				scoopStatusOutput = createScoopStatusOutput([]*Package{
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
				})
			})

			It("should still return results in order of package name", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				}))
			})
		})

		When("`scoop status` doesn't include any packages", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = "Scoop is up to date.\nEverything is ok!"
			})

			It("should indicate that all packages are up to date", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.0"),
					NewPackageFromStrings("package2", "2.3.4", "2.3.4"),
					NewPackageFromStrings("package3", "3.2.1", "3.2.1"),
				}))
			})
		})

		When("`scoop status` includes some but not all packages", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = createScoopStatusOutput([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				})
			})

			It("should indicate that some but not all packages are up to date", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.1"),
					NewPackageFromStrings("package2", "2.3.4", "2.3.4"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				}))
			})
		})

		When("scoop has no packages installed", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{})
				scoopStatusOutput = createScoopStatusOutput([]*Package{})
			})

			It("should return an empty slice", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{}))
			})
		})

		When("`scoop export` returns invalid JSON", func() {
			BeforeEach(func() {
				scoopExportOutput = "invalid json"
				scoopStatusOutput = createScoopStatusOutput([]*Package{})
			})

			It("should return a `json.SyntaxError`", func() {
				_, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeAssignableToTypeOf(&json.SyntaxError{}))
			})
		})

		When("`scoop status` output is just arbitrary text with no results", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = "arbitrary text"
			})

			It("should assume all packages are up to date", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.0"),
					NewPackageFromStrings("package2", "2.3.4", "2.3.4"),
					NewPackageFromStrings("package3", "3.2.1", "3.2.1"),
				}))
			})
		})

		When("a result in `scoop status` has the wrong number of fields", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = `Scoop is up to date.

Name Installed Version Latest Version Missing Dependencies Info
---- ----------------- -------------- -------------------- ----
package1 1.0.0 1.0.1
package2
package3 3.2.1 4.0.0
`
			})

			It("should return an error", func() {
				_, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(Not(BeNil()))
			})
		})
	})

	Describe("InstallPackage", func() {
	})

	Describe("UpdatePackage", func() {
	})

	Describe("UninstallPackage", func() {
	})
})

type scoopExportData struct {
	Buckets []scoopExportBucket `json:"buckets"`
	Apps    []scoopExportApp    `json:"apps"`
}

type scoopExportBucket struct {
	Name string `json:"name"`
}

type scoopExportApp struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
}

// createScoopExportOutput creates `scoop export` output from the given information.
func createScoopExportOutput(apps []scoopExportApp) string {
	bucketNames := make(map[string]bool)
	for _, app := range apps {
		bucketNames[app.Source] = true
	}

	var buckets []scoopExportBucket
	for bucketName := range bucketNames {
		buckets = append(buckets, scoopExportBucket{bucketName})
	}

	data := scoopExportData{
		Buckets: buckets,
		Apps:    apps,
	}

	if bytes, err := json.Marshal(data); err == nil {
		return string(bytes)
	}
	return ""
}

// createScoopStatusOutput creates `scoop status` output from the given information.
func createScoopStatusOutput(packages []*Package) string {
	var outputStringBuilder strings.Builder
	outputStringBuilder.WriteString("Scoop is up to date.\n\n")
	outputStringBuilder.WriteString("Name Installed Version Latest Version Missing Dependencies Info\n")
	outputStringBuilder.WriteString("---- ----------------- -------------- -------------------- ----\n")

	for _, pkg := range packages {
		outputStringBuilder.WriteString(fmt.Sprintf("%s %s %s\n", pkg.Name, pkg.InstalledVersion, pkg.LatestVersion))
	}

	return outputStringBuilder.String()
}
