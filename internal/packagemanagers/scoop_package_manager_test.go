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
		It("should return true if running on Windows", func() {
			operatingSystemServiceDouble.SetIsWindows(true)
			result := scoopPackageManager.IsSupported()
			Expect(result).To(BeTrue())
		})

		It("should return false if not running on Windows", func() {
			operatingSystemServiceDouble.SetIsWindows(false)
			result := scoopPackageManager.IsSupported()
			Expect(result).To(BeFalse())
		})
	})

	Describe("IsInstalled", func() {
		It("should write output stating it is checking if Scoop is installed", func() {
			scoopPackageManager.IsInstalled()
			Expect(outputWriterDouble.String()).To(Equal("Checking if package manager \"scoop\" is installed...\n"))
		})

		It("should return true if `scoop --version` runs successfully", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("scoop output", "scoop", false, "--version")
			result, err := scoopPackageManager.IsInstalled()
			Expect(err).To(BeNil())
			Expect(result).To(BeTrue())
		})

		It("should return false if running `scoop --version` returns an error", func() {
			result, err := scoopPackageManager.IsInstalled()
			Expect(err).To(BeNil())
			Expect(result).To(BeFalse())
		})
	})

	Describe("Install", func() {
		It("should write output stating it is installing Scoop", func() {
			scoopPackageManager.Install()
			Expect(outputWriterDouble.String()).To(Equal("Installing package manager \"scoop\"...\n"))
		})

		It("should run the Scoop install process and show its output", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("scoop installation output", "powershell", true,
				"irm get.scoop.sh | iex")
			err := scoopPackageManager.Install()
			Expect(err).To(BeNil())
			Expect(shellCommandServiceDouble.WasCalledWith("powershell", true, "irm get.scoop.sh | iex")).To(BeTrue())
		})

		It("should return an error if the Scoop install process returns an error", func() {
			err := scoopPackageManager.Install()
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("should write output stating it is updating Scoop", func() {
			scoopPackageManager.Update()
			Expect(outputWriterDouble.String()).To(Equal("Updating package manager \"scoop\"...\n"))
		})

		It("should run `scoop update` and show its output", func() {
			commandOutput := `some output
some more output
Scoop was updated successfully!`
			shellCommandServiceDouble.SetOutputForExpectedInputs(commandOutput, "scoop", true, "update")
			err := scoopPackageManager.Update()
			Expect(err).To(BeNil())
			Expect(shellCommandServiceDouble.WasCalledWith("scoop", true, "update")).To(BeTrue())
		})

		It("should return an error if `scoop update` returns an error", func() {
			err := scoopPackageManager.Update()
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `scoop update` output does not contain \"Scoop was updated\"", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("scoop update output", "scoop", true, "update")
			err := scoopPackageManager.Update()
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Uninstall", func() {
		It("should write output stating it is uninstalling Scoop", func() {
			scoopPackageManager.Uninstall()
			Expect(outputWriterDouble.String()).To(Equal("Uninstalling package manager \"scoop\"...\n"))
		})

		It("should run `scoop uninstall scoop` and show its output", func() {
			commandOutput := `some output
some more output
'scoop' was uninstalled.`
			shellCommandServiceDouble.SetOutputForExpectedInputs(commandOutput, "scoop", true, "uninstall", "scoop")
			err := scoopPackageManager.Uninstall()
			Expect(err).To(BeNil())
			Expect(shellCommandServiceDouble.WasCalledWith("scoop", true, "uninstall", "scoop")).To(BeTrue())
		})

		It("should return an error if `scoop uninstall scoop` returns an error", func() {
			err := scoopPackageManager.Uninstall()
			Expect(err).NotTo(BeNil())
		})

		It("should return an error if `scoop uninstall scoop` output does not contain \"'scoop' was uninstalled\"",
			func() {
				shellCommandServiceDouble.SetOutputForExpectedInputs("scoop uninstall output", "scoop", true,
					"uninstall", "scoop")
				err := scoopPackageManager.Uninstall()
				Expect(err).NotTo(BeNil())
			})
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

		When("Scoop has no packages installed", func() {
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
		const packageName = "package1"
		const packageVersion = "1.0.0"

		It("should write output stating it is installing the given package", func() {
			scoopPackageManager.InstallPackage(packageName, nil)
			Expect(outputWriterDouble.String()).To(Equal(fmt.Sprintf("Installing package \"%s\"...\n", packageName)))
		})

		It("should capture the version number from the `scoop install` command's output and return it", func() {
			commandOutput := `some output
some more output
'%s' (%s) was installed successfully!`
			commandOutput = fmt.Sprintf(commandOutput, packageName, packageVersion)
			shellCommandServiceDouble.SetOutputForExpectedInputs(commandOutput, "scoop", true, "install", packageName)
			result, err := scoopPackageManager.InstallPackage(packageName, nil)
			Expect(err).To(BeNil())
			Expect(result.String()).To(Equal(packageVersion))
		})

		It("should return an error if the `scoop install` command returns an error", func() {
			_, err := scoopPackageManager.InstallPackage(packageName, nil)
			Expect(err).To(Not(BeNil()))
		})

		It("should return an error if the `scoop install` command output does not contain \"'<package>' (<version>) "+
			"was installed\"", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("some output", "scoop", true, "install", packageName)
			_, err := scoopPackageManager.InstallPackage(packageName, nil)
			Expect(err).To(Not(BeNil()))
		})
	})

	Describe("UpdatePackage", func() {
		const packageName = "package1"
		const packageVersion = "1.0.1"

		It("should write output stating it is updating the given package", func() {
			scoopPackageManager.UpdatePackage(packageName, nil)
			Expect(outputWriterDouble.String()).To(Equal(fmt.Sprintf("Updating package \"%s\"...\n", packageName)))
		})

		It("should capture the version number from the `scoop update` command's output and return it", func() {
			commandOutput := `some output
some more output
'%s' (%s) was installed successfully!`
			commandOutput = fmt.Sprintf(commandOutput, packageName, packageVersion)
			shellCommandServiceDouble.SetOutputForExpectedInputs(commandOutput, "scoop", true, "update", packageName)
			result, err := scoopPackageManager.UpdatePackage(packageName, nil)
			Expect(err).To(BeNil())
			Expect(result.String()).To(Equal(packageVersion))
		})

		It("should return an error if the `scoop update` command returns an error", func() {
			_, err := scoopPackageManager.UpdatePackage(packageName, nil)
			Expect(err).To(Not(BeNil()))
		})

		It("should return an error if the `scoop update` command output does not contain \"'<package>' (<version>) "+
			"was installed\"", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("some output", "scoop", true, "update", packageName)
			_, err := scoopPackageManager.UpdatePackage(packageName, nil)
			Expect(err).To(Not(BeNil()))
		})
	})

	Describe("UninstallPackage", func() {
		const packageName = "package1"

		It("should write output stating it is uninstalling the given package", func() {
			scoopPackageManager.UninstallPackage(packageName)
			Expect(outputWriterDouble.String()).To(Equal(fmt.Sprintf("Uninstalling package \"%s\"...\n", packageName)))
		})

		It("should run the `scoop uninstall` command and show its output", func() {
			commandOutput := `some output
some more output
'%s' was uninstalled.`
			commandOutput = fmt.Sprintf(commandOutput, packageName)
			shellCommandServiceDouble.SetOutputForExpectedInputs(commandOutput, "scoop", true, "uninstall", packageName)
			err := scoopPackageManager.UninstallPackage(packageName)
			Expect(err).To(BeNil())
		})

		It("should return an error if the `scoop uninstall` command returns an error", func() {
			err := scoopPackageManager.UninstallPackage(packageName)
			Expect(err).To(Not(BeNil()))
		})

		It("should return an error if the `scoop uninstall` command output does not contain \"'<package>' "+
			"was uninstalled\"", func() {
			shellCommandServiceDouble.SetOutputForExpectedInputs("some output", "scoop", true, "uninstall", packageName)
			err := scoopPackageManager.UninstallPackage(packageName)
			Expect(err).To(Not(BeNil()))
		})
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
