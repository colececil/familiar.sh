package packagemanagers_test

import (
	. "github.com/colececil/familiar.sh/internal/packagemanagers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
)

var _ = Describe("PackageManagerRegistry", func() {
	const scoop = "scoop"
	const chocolatey = "chocolatey"
	const homebrew = "homebrew"

	var packageManagerRegistry PackageManagerRegistry
	var scoopPackageManager PackageManager
	var chocolateyPackageManager PackageManager
	var homebrewPackageManager PackageManager

	BeforeEach(func() {
		mock.SetUp(GinkgoT())

		scoopPackageManager = mock.Mock[PackageManager]()
		mock.WhenSingle(scoopPackageManager.Name()).ThenReturn(scoop)
		mock.WhenSingle(scoopPackageManager.Order()).ThenReturn(1)

		chocolateyPackageManager = mock.Mock[PackageManager]()
		mock.WhenSingle(chocolateyPackageManager.Name()).ThenReturn(chocolatey)
		mock.WhenSingle(chocolateyPackageManager.Order()).ThenReturn(2)

		homebrewPackageManager = mock.Mock[PackageManager]()
		mock.WhenSingle(homebrewPackageManager.Name()).ThenReturn(homebrew)
		mock.WhenSingle(homebrewPackageManager.Order()).ThenReturn(3)

		packageManagerRegistry = NewPackageManagerRegistry([]PackageManager{
			homebrewPackageManager,
			chocolateyPackageManager,
			scoopPackageManager,
		})
	})

	Describe("NewPackageManagerRegistry", func() {
		const panicMessage = "package manager registry does not contain the expected package managers"

		It("should panic if the package manager registry does not contain the expected package managers", func() {
			Expect(func() {
				NewPackageManagerRegistry([]PackageManager{
					scoopPackageManager,
					chocolateyPackageManager,
				})
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the package managers' orders are less than 1", func() {
			mock.SetUp(GinkgoT())

			scoopPackageManager = mock.Mock[PackageManager]()
			mock.WhenSingle(scoopPackageManager.Name()).ThenReturn(scoop)
			mock.WhenSingle(scoopPackageManager.Order()).ThenReturn(0)

			Expect(func() {
				NewPackageManagerRegistry([]PackageManager{
					scoopPackageManager,
					chocolateyPackageManager,
					homebrewPackageManager,
				})
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the package managers' orders are greater than the number of package managers",
			func() {
				mock.SetUp(GinkgoT())

				homebrewPackageManager = mock.Mock[PackageManager]()
				mock.WhenSingle(homebrewPackageManager.Name()).ThenReturn(homebrew)
				mock.WhenSingle(homebrewPackageManager.Order()).ThenReturn(4)

				Expect(func() {
					NewPackageManagerRegistry([]PackageManager{
						scoopPackageManager,
						chocolateyPackageManager,
						homebrewPackageManager,
					})
				}).To(PanicWith(panicMessage))
			})

		It("should panic if any of the package managers' orders are not unique", func() {
			Expect(func() {
				mock.SetUp(GinkgoT())

				chocolateyPackageManager = mock.Mock[PackageManager]()
				mock.WhenSingle(chocolateyPackageManager.Name()).ThenReturn(chocolatey)
				mock.WhenSingle(chocolateyPackageManager.Order()).ThenReturn(1)

				NewPackageManagerRegistry([]PackageManager{
					scoopPackageManager,
					chocolateyPackageManager,
					homebrewPackageManager,
				})
			}).To(PanicWith(panicMessage))
		})
	})

	Describe("GetAllPackageManagers", func() {
		It("should return a slice containing all package managers, sorted by their returned orders", func() {
			result := packageManagerRegistry.GetAllPackageManagers()
			Expect(len(result)).To(Equal(3))
			Expect(result[0]).To(Equal(scoopPackageManager))
			Expect(result[1]).To(Equal(chocolateyPackageManager))
			Expect(result[2]).To(Equal(homebrewPackageManager))
		})
	})

	Describe("GetPackageManager", func() {
		It("should return the package manager of the given name if it is in the registry", func() {
			result, err := packageManagerRegistry.GetPackageManager(chocolatey)
			Expect(err).To(BeNil())
			Expect(result).To(Equal(chocolateyPackageManager))
		})

		It("should return an error if the package manager with the given name is not in the registry", func() {
			_, err := packageManagerRegistry.GetPackageManager("invalid")
			Expect(err.Error()).To(Equal("package manager not valid"))
		})
	})
})
