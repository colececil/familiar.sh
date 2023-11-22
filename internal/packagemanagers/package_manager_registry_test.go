package packagemanagers_test

import (
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PackageManagerRegistry", func() {
	const packageManager1Name = "packageManager1"
	const packageManager2Name = "packageManager2"
	const packageManager3Name = "packageManager3"

	var packageManagerRegistry PackageManagerRegistry
	var packageManager1 PackageManager
	var packageManager2 PackageManager
	var packageManager3 PackageManager

	BeforeEach(func() {
		packageManager1 = test.NewPackageManagerDouble(packageManager1Name)
		packageManager2 = test.NewPackageManagerDouble(packageManager2Name)
		packageManager3 = test.NewPackageManagerDouble(packageManager3Name)
		packageManagerRegistry = NewPackageManagerRegistry([]PackageManager{
			packageManager1,
			packageManager2,
			packageManager3,
		})
	})

	Describe("GetAllPackageManagers", func() {
		It("should return a slice containing all package managers", func() {
			result := packageManagerRegistry.GetAllPackageManagers()
			Expect(len(result)).To(Equal(3))
			Expect(result[0]).To(Equal(packageManager1))
			Expect(result[1]).To(Equal(packageManager2))
			Expect(result[2]).To(Equal(packageManager3))
		})
	})

	Describe("GetPackageManager", func() {
		It("should return the package manager of the given name if it is in the registry", func() {
			result, err := packageManagerRegistry.GetPackageManager("packageManager2")
			Expect(err).To(BeNil())
			Expect(result).To(Equal(packageManager2))
		})

		It("should return an error if the package manager with the given name is not in the registry", func() {
			_, err := packageManagerRegistry.GetPackageManager("invalid")
			Expect(err.Error()).To(Equal("package manager not valid"))
		})
	})
})
