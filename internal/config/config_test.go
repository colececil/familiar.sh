package config_test

import (
	. "github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Config", func() {
	const packageManager1Name = "packageManager1"
	const packageManager2Name = "packageManager2"
	const packageManager3Name = "packageManager3"
	const package1Name = "package1"
	const package2Name = "package2"
	const package3Name = "package3"

	var config *Config
	var packageManagerRegistry packagemanagers.PackageManagerRegistry

	BeforeEach(func() {
		config = NewConfig()
		packageManagerRegistry = packagemanagers.NewPackageManagerRegistry([]packagemanagers.PackageManager{
			test.NewPackageManagerDouble(packageManager1Name),
			test.NewPackageManagerDouble(packageManager2Name),
			test.NewPackageManagerDouble(packageManager3Name),
		})
	})

	Describe("AddPackageManager", func() {
		It("should add the given package manager to the config and keep the list in alphabetical order", func() {
			err := config.AddPackageManager(packageManager3Name, packageManagerRegistry)
			Expect(err).To(BeNil())
			Expect(len(config.PackageManagers)).To(Equal(1))
			Expect(config.PackageManagers[0].Name).To(Equal(packageManager3Name))

			err = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
			Expect(err).To(BeNil())
			Expect(len(config.PackageManagers)).To(Equal(2))
			Expect(config.PackageManagers[0].Name).To(Equal(packageManager2Name))
			Expect(config.PackageManagers[1].Name).To(Equal(packageManager3Name))

			err = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			Expect(err).To(BeNil())
			Expect(len(config.PackageManagers)).To(Equal(3))
			Expect(config.PackageManagers[0].Name).To(Equal(packageManager1Name))
			Expect(config.PackageManagers[1].Name).To(Equal(packageManager2Name))
			Expect(config.PackageManagers[2].Name).To(Equal(packageManager3Name))
		})

		It("should return an error if the given package manager is not in the package manager registry", func() {
			err := config.AddPackageManager("invalid", packageManagerRegistry)
			Expect(err.Error()).To(Equal("package manager not valid"))
			Expect(len(config.PackageManagers)).To(Equal(0))
		})

		It("should return an error if the given package manager is already in the config", func() {
			err := config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			Expect(err).To(BeNil())
			Expect(len(config.PackageManagers)).To(Equal(1))

			err = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			Expect(err.Error()).To(Equal("package manager already present"))
			Expect(len(config.PackageManagers)).To(Equal(1))
		})
	})

	Describe("RemovePackageManager", func() {
		BeforeEach(func() {
			_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
			_ = config.AddPackageManager(packageManager3Name, packageManagerRegistry)
		})

		It("should remove the given package manager from the config", func() {
			err := config.RemovePackageManager(packageManager2Name)
			Expect(err).To(BeNil())
			Expect(len(config.PackageManagers)).To(Equal(2))
			Expect(config.PackageManagers[0].Name).To(Equal(packageManager1Name))
			Expect(config.PackageManagers[1].Name).To(Equal(packageManager3Name))
		})

		It("should return an error if the given package manager is not in the config", func() {
			err := config.RemovePackageManager("invalid")
			Expect(err.Error()).To(Equal("package manager not present"))
			Expect(len(config.PackageManagers)).To(Equal(3))
			Expect(config.PackageManagers[0].Name).To(Equal(packageManager1Name))
			Expect(config.PackageManagers[1].Name).To(Equal(packageManager2Name))
			Expect(config.PackageManagers[2].Name).To(Equal(packageManager3Name))
		})
	})

	Describe("AddPackage", func() {
		var packageManager1 *ConfiguredPackageManager
		var packageManager2 *ConfiguredPackageManager

		BeforeEach(func() {
			_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
			packageManager1 = &config.PackageManagers[0]
			packageManager2 = &config.PackageManagers[1]
		})

		It("should add the given version of the given package to the config under the given package manager, keeping "+
			"the list in alphabetical order", func() {

			package1Version := packagemanagers.NewVersion("1.0.0")
			package2Version := packagemanagers.NewVersion("1.1.1")
			package3Version := packagemanagers.NewVersion("1.0.1")

			err := config.AddPackage(packageManager2Name, package3Name, package3Version)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(0))
			Expect(len(packageManager2.Packages)).To(Equal(1))
			Expect(packageManager2.Packages[0].Name).To(Equal(package3Name))
			Expect(packageManager2.Packages[0].Version).To(Equal(package3Version.VersionString))

			err = config.AddPackage(packageManager2Name, package2Name, package2Version)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(0))
			Expect(len(packageManager2.Packages)).To(Equal(2))
			Expect(packageManager2.Packages[0].Name).To(Equal(package2Name))
			Expect(packageManager2.Packages[0].Version).To(Equal(package2Version.VersionString))
			Expect(packageManager2.Packages[1].Name).To(Equal(package3Name))
			Expect(packageManager2.Packages[1].Version).To(Equal(package3Version.VersionString))

			err = config.AddPackage(packageManager2Name, package1Name, package1Version)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(0))
			Expect(len(packageManager2.Packages)).To(Equal(3))
			Expect(packageManager2.Packages[0].Name).To(Equal(package1Name))
			Expect(packageManager2.Packages[0].Version).To(Equal(package1Version.VersionString))
			Expect(packageManager2.Packages[1].Name).To(Equal(package2Name))
			Expect(packageManager2.Packages[1].Version).To(Equal(package2Version.VersionString))
			Expect(packageManager2.Packages[2].Name).To(Equal(package3Name))
			Expect(packageManager2.Packages[2].Version).To(Equal(package3Version.VersionString))
		})

		It("should return an error if the given package manager is not in the config", func() {
			err := config.AddPackage("invalid", package1Name, packagemanagers.NewVersion("1.0.0"))
			Expect(err.Error()).To(Equal("package manager not present"))
			Expect(len(packageManager1.Packages)).To(Equal(0))
			Expect(len(packageManager2.Packages)).To(Equal(0))
		})

		It("should return an error if the given package is already in the config under the given package manager",
			func() {
				err := config.AddPackage(packageManager1Name, package1Name, packagemanagers.NewVersion("1.0.0"))
				Expect(err).To(BeNil())
				Expect(len(packageManager1.Packages)).To(Equal(1))

				err = config.AddPackage(packageManager1Name, package1Name, packagemanagers.NewVersion("1.0.1"))
				Expect(err.Error()).To(Equal("package already present"))
				Expect(len(packageManager1.Packages)).To(Equal(1))
			})

		It("should not return an error if the given package is already in the config, but only under a different "+
			"package manager", func() {

			packageVersion := packagemanagers.NewVersion("1.0.0")

			err := config.AddPackage(packageManager1Name, package1Name, packageVersion)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(1))
			Expect(len(packageManager2.Packages)).To(Equal(0))

			err = config.AddPackage(packageManager2Name, package1Name, packageVersion)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(1))
			Expect(len(packageManager2.Packages)).To(Equal(1))
		})
	})

	Describe("UpdatePackage", func() {
		const packageStartingVersionString = "1.0.0"
		const packageUpdateVersionString = "1.0.1"

		var packageManager1 *ConfiguredPackageManager
		var packageManager2 *ConfiguredPackageManager

		BeforeEach(func() {
			_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			_ = config.AddPackage(packageManager1Name, package1Name,
				packagemanagers.NewVersion(packageStartingVersionString))
			_ = config.AddPackage(packageManager1Name, package2Name,
				packagemanagers.NewVersion(packageStartingVersionString))

			_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
			_ = config.AddPackage(packageManager2Name, package1Name,
				packagemanagers.NewVersion(packageStartingVersionString))
			_ = config.AddPackage(packageManager2Name, package2Name,
				packagemanagers.NewVersion(packageStartingVersionString))

			packageManager1 = &config.PackageManagers[0]
			packageManager2 = &config.PackageManagers[1]
		})

		It("should update the given package to the given version under the given package manager", func() {
			err := config.UpdatePackage(packageManager1Name, package1Name,
				packagemanagers.NewVersion(packageUpdateVersionString))
			Expect(err).To(BeNil())
			Expect(packageManager1.Packages[0].Version).To(Equal(packageUpdateVersionString))
			Expect(packageManager1.Packages[1].Version).To(Equal(packageStartingVersionString))
			Expect(packageManager2.Packages[0].Version).To(Equal(packageStartingVersionString))
			Expect(packageManager2.Packages[1].Version).To(Equal(packageStartingVersionString))

			err = config.UpdatePackage(packageManager2Name, package1Name,
				packagemanagers.NewVersion(packageUpdateVersionString))
			Expect(err).To(BeNil())
			Expect(packageManager1.Packages[0].Version).To(Equal(packageUpdateVersionString))
			Expect(packageManager1.Packages[1].Version).To(Equal(packageStartingVersionString))
			Expect(packageManager2.Packages[0].Version).To(Equal(packageUpdateVersionString))
			Expect(packageManager2.Packages[1].Version).To(Equal(packageStartingVersionString))
		})

		It("should return an error if the given package manager is not in the config", func() {
			err := config.UpdatePackage("invalid", package1Name,
				packagemanagers.NewVersion(packageUpdateVersionString))
			Expect(err.Error()).To(Equal("package manager not present"))
		})

		It("should return an error if the given package is not in the config under the given package manager", func() {
			err := config.UpdatePackage(packageManager1Name, "invalid",
				packagemanagers.NewVersion(packageUpdateVersionString))
			Expect(err.Error()).To(Equal("package not present"))
		})

		It("should return an error if the given package is already set to the given version under the given package "+
			"manager", func() {
			err := config.UpdatePackage(packageManager1Name, package1Name,
				packagemanagers.NewVersion(packageStartingVersionString))
			Expect(err.Error()).To(Equal("package already set to given version"))
		})
	})

	Describe("RemovePackage", func() {
		var packageManager1 *ConfiguredPackageManager
		var packageManager2 *ConfiguredPackageManager

		BeforeEach(func() {
			const packageVersionString = "1.0.0"
			_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
			_ = config.AddPackage(packageManager1Name, package1Name, packagemanagers.NewVersion(packageVersionString))
			_ = config.AddPackage(packageManager1Name, package2Name, packagemanagers.NewVersion(packageVersionString))

			_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
			_ = config.AddPackage(packageManager2Name, package3Name, packagemanagers.NewVersion(packageVersionString))

			packageManager1 = &config.PackageManagers[0]
			packageManager2 = &config.PackageManagers[1]
		})

		It("should remove the given package from the config under the given package manager", func() {
			err := config.RemovePackage(packageManager1Name, package1Name)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(1))
			Expect(packageManager1.Packages[0].Name).To(Equal(package2Name))

			err = config.RemovePackage(packageManager1Name, package2Name)
			Expect(err).To(BeNil())
			Expect(len(packageManager1.Packages)).To(Equal(0))
		})

		It("should return an error if the given package manager is not in the config", func() {
			err := config.RemovePackage("invalid", package1Name)
			Expect(err.Error()).To(Equal("package manager not present"))
			Expect(len(packageManager1.Packages)).To(Equal(2))
			Expect(len(packageManager2.Packages)).To(Equal(1))
		})

		It("should return an error if the given package is not in the config under the given package manager", func() {
			err := config.RemovePackage(packageManager2Name, package1Name)
			Expect(err.Error()).To(Equal("package not present"))
			Expect(len(packageManager1.Packages)).To(Equal(2))
			Expect(len(packageManager2.Packages)).To(Equal(1))
		})
	})

	Describe("YamlString", func() {
		DescribeTable("should return the correct YAML representation of the config",
			func(expectedYamlString string, configSetupFunc func()) {
				configSetupFunc()
				result, err := config.YamlString()
				Expect(err).To(BeNil())

				parsedExpected, err := parseYamlString(expectedYamlString)
				Expect(err).To(BeNil())
				parsedResult, err := parseYamlString(result)
				Expect(err).To(BeNil())

				Expect(parsedResult).To(Equal(parsedExpected))
			},
			Entry("when the config is empty",
				`version: 1
files: []
scripts: []
packageManagers: []
`,
				func() {},
			),
			Entry("when the config contains package managers with no packages",
				`version: 1
files: []
scripts: []
packageManagers:
  - name: packageManager1
    packages: []
  - name: packageManager2
    packages: []
`,
				func() {
					_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
					_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)
				},
			),
			Entry("when the config contains package managers with packages",
				`version: 1
files: []
scripts: []
packageManagers:
  - name: packageManager1
    packages:
      - name: package1
        version: 1.0.0
      - name: package2
        version: 1.1.1
  - name: packageManager2
    packages: []
  - name: packageManager3
    packages:
      - name: package3
        version: 1.0.1
`,
				func() {
					_ = config.AddPackageManager(packageManager1Name, packageManagerRegistry)
					_ = config.AddPackage(packageManager1Name, package1Name, packagemanagers.NewVersion("1.0.0"))
					_ = config.AddPackage(packageManager1Name, package2Name, packagemanagers.NewVersion("1.1.1"))

					_ = config.AddPackageManager(packageManager2Name, packageManagerRegistry)

					_ = config.AddPackageManager(packageManager3Name, packageManagerRegistry)
					_ = config.AddPackage(packageManager3Name, package3Name, packagemanagers.NewVersion("1.0.1"))
				},
			),
		)
	})
})

// parseYamlString parses the given YAML string and returns the result.
func parseYamlString(yamlString string) (any, error) {
	var result any
	err := yaml.Unmarshal([]byte(yamlString), &result)
	return result, err
}
