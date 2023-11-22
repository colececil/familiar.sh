package config_test

import (
	. "github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Config", func() {
	var config *Config
	var packageManagerRegistry packagemanagers.PackageManagerRegistry

	BeforeEach(func() {
		config = NewConfig()
		packageManagerRegistry = packagemanagers.NewPackageManagerRegistry([]packagemanagers.PackageManager{})
	})

	Describe("AddPackageManager", func() {
		It("should append the given package manager to the end of the config's list of package managers", func() {
		})

		It("should return an error if the given package manager is not in the package manager registry", func() {
		})

		It("should return an error if the given package manager is already in the config", func() {
		})
	})
})
