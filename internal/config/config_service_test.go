package config_test

import (
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/system"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigService", func() {
	const configLocation = "/path/to/config"

	var fileSystemServiceDouble system.FileSystemService
	var configService *config.ConfigService

	BeforeEach(func() {
		configService = config.NewConfigService(fileSystemServiceDouble)
	})

	Describe("GetConfigLocation", func() {
		It("should return the config location defined in the config location file", func() {
			location, err := configService.GetConfigLocation()

			Expect(err).To(BeNil())
			Expect(location).To(Equal(configLocation))
		})

		It("should return an empty string if the config location file does not exist", func() {
		})

		It("should return an empty string if the config location file is empty", func() {
		})
	})

	Describe("SetConfigLocation", func() {
	})

	Describe("GetConfig", func() {
	})

	Describe("SetConfig", func() {
	})
})
