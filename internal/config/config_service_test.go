package config_test

import (
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigService", func() {
	const configHomeLocation = "/home/user/.config"
	const configLocation = "/path/to/config"

	var fileSystemServiceDouble *test.FileSystemServiceDouble
	var configService *config.ConfigService

	BeforeEach(func() {
		fileSystemServiceDouble = test.NewFileSystemServiceDouble()
		fileSystemServiceDouble.SetXdgConfigHome(configHomeLocation)
		fileSystemServiceDouble.SetFileContentForExpectedPath(
			configHomeLocation+"/io.colececil.familiar/config_location",
			configLocation,
		)
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
