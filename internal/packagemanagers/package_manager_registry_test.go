package packagemanagers_test

import (
	"bytes"
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/system"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PackageManagerRegistry", func() {
	var operatingSystemService *system.OperatingSystemService
	var shellCommandServiceDouble *test.ShellCommandServiceDouble
	var outputWriterDouble *bytes.Buffer
	var scoopPackageManager *ScoopPackageManager
	var packageManagerRegistry PackageManagerRegistry

	BeforeEach(func() {
		operatingSystemService = system.NewOperatingSystemService(system.WindowsOperatingSystem)
		shellCommandServiceDouble = test.NewShellCommandServiceDouble()
		outputWriterDouble = new(bytes.Buffer)
		scoopPackageManager = NewScoopPackageManager(operatingSystemService,
			shellCommandServiceDouble.ShellCommandService, outputWriterDouble)
		packageManagerRegistry = NewPackageManagerRegistry(scoopPackageManager)
	})

	Describe("GetAllPackageManagers", func() {
		It("should return a slice containing all package managers", func() {
			result := packageManagerRegistry.GetAllPackageManagers()
			Expect(len(result)).To(Equal(1))
			Expect(result[0]).To(Equal(scoopPackageManager))
		})
	})

	Describe("GetPackageManager", func() {
		It("should return the correct package manager", func() {
			result, err := packageManagerRegistry.GetPackageManager("scoop")
			Expect(err).To(BeNil())
			Expect(result).To(Equal(scoopPackageManager))
		})

		It("should return an error if the package manager with the given name is not in the registry", func() {
			_, err := packageManagerRegistry.GetPackageManager("invalid")
			Expect(err.Error()).To(Equal("package manager not valid"))
		})
	})
})
