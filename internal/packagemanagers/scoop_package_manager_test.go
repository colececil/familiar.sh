package packagemanagers_test

import (
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

		JustBeforeEach(func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopExportOutput, "scoop", false, "export")
			shellCommandServiceDouble.SetOutputForExpectedInputs(scoopStatusOutput, "scoop", false, "status")
		})

		When("all packages returned by `scoop export` are included in `scoop status`", func() {
			BeforeEach(func() {
				scoopExportOutput = createScoopExportOutput([]scoopExportApp{
					{Name: "package1", Version: "1.0.0", Source: "main"},
					{Name: "package2", Version: "2.3.4", Source: "main"},
					{Name: "package3", Version: "3.2.1", Source: "main"},
				})
				scoopStatusOutput = createScoopStatusOutput([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.0"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				})
			})

			It("should indicate that all packages have updates available", func() {
				packages, err := scoopPackageManager.InstalledPackages()
				Expect(err).To(BeNil())
				Expect(packages).To(Equal([]*Package{
					NewPackageFromStrings("package1", "1.0.0", "1.0.0"),
					NewPackageFromStrings("package2", "2.3.4", "2.5.0"),
					NewPackageFromStrings("package3", "3.2.1", "4.0.0"),
				}))
			})
		})

		PIt("should sort the packages by name", func() {
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

	buckets := make([]scoopExportBucket, 0)
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
