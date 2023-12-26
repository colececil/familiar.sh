package config_test

import (
	"bytes"
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("ConfigService", func() {
	const configHomeLocation = "/home/user/.config"
	const appDirectoryName = "io.colececil.familiar"
	const configLocationFileName = "config_location"
	const configLocation = "/path/to/config.yml"

	var fileSystemServiceDouble *test.FileSystemServiceDouble
	var outputWriterDouble *bytes.Buffer
	var configService *config.ConfigService

	BeforeEach(func() {
		fileSystemServiceDouble = test.NewFileSystemServiceDouble()
		fileSystemServiceDouble.SetXdgConfigHome(configHomeLocation)
		outputWriterDouble = new(bytes.Buffer)
		configService = config.NewConfigService(fileSystemServiceDouble, outputWriterDouble)
	})

	Describe("GetConfigLocation", func() {
		It("should return the config location defined in the config location file", func() {
			file, _ := fileSystemServiceDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
				configLocationFileName)
			_, _ = file.Write([]byte(configLocation))
			_ = file.Close()

			location, err := configService.GetConfigLocation()

			Expect(err).To(BeNil())
			Expect(location).To(Equal(configLocation))
		})

		It("should return an error if the config location file does not exist", func() {
			expectedError := fmt.Errorf("The location of Familiar's shared config file has not yet been set. Please " +
				"set it using \"familiar config location <path>\", for more details, execute \"familiar help config\".")

			location, err := configService.GetConfigLocation()

			Expect(err).To(Equal(expectedError))
			Expect(location).To(Equal(""))
		})

		It("should return an error if the config location file is empty", func() {
			expectedError := fmt.Errorf("The location of Familiar's shared config file has not yet been set. Please " +
				"set it using \"familiar config location <path>\", for more details, execute \"familiar help config\".")
			file, _ := fileSystemServiceDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
				configLocationFileName)
			_ = file.Close()

			location, err := configService.GetConfigLocation()

			Expect(err).To(Equal(expectedError))
			Expect(location).To(Equal(""))
		})
	})

	Describe("SetConfigLocation", func() {
		When("the config location file exists", func() {
			var expectedFileContent string

			BeforeEach(func() {
				file, _ := fileSystemServiceDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
					configLocationFileName)
				_ = file.Close()

				configFileDir := fileSystemServiceDouble.Dir(configLocation)
				file, _ = fileSystemServiceDouble.CreateFile(configFileDir)
				_ = file.Close()

				expectedFileContent, _ = fileSystemServiceDouble.Abs(configLocation)
				expectedFileContent += "\n"
			})

			It("should write the given path to the config location file", func() {
				err := configService.SetConfigLocation(configLocation)
				fileContentBytes, _ := fileSystemServiceDouble.ReadFile(configHomeLocation + "/" +
					appDirectoryName + "/" + configLocationFileName)
				fileContent := string(fileContentBytes)

				Expect(err).To(BeNil())
				Expect(fileContent).To(Equal(expectedFileContent))
			})

			It("should write output stating that the config location has been set", func() {
				err := configService.SetConfigLocation(configLocation)
				Expect(err).To(BeNil())
				Expect(outputWriterDouble.String()).To(Equal("The config file location has been set to \"" +
					configLocation + "\".\n"))
			})

			It("should close the config location file after writing to it", func() {
				err := configService.SetConfigLocation(configLocation)
				file, _ := fileSystemServiceDouble.GetCreatedFile(configHomeLocation + "/" + appDirectoryName + "/" +
					configLocationFileName)

				Expect(err).To(BeNil())
				Expect(file.IsClosed()).To(BeTrue())
			})

			It("should return an error if the directory of the file specified by the given path does not exist", func() {
				dir := "/other/dir/to"
				err := configService.SetConfigLocation(dir + "/config.yml")
				Expect(err.Error()).To(Equal(fmt.Sprintf("directory \"%s\" does not exist", dir)))
			})

			It("should return an error if the file specified by the path does not have a YAML file extension", func() {
				path := strings.TrimSuffix(configLocation, ".yml")
				err := configService.SetConfigLocation(path)
				Expect(err.Error()).To(Equal(fmt.Sprintf(
					"invalid file extension \"\": expected \".yml\" or \".yaml\"")))
			})

			It("should return an error if there is an issue checking if the directory exists", func() {
			})

			It("should return an error if there is an issue updating the config location file", func() {
			})
		})

		When("the config location file does not exist", func() {
			It("should create the config location file", func() {
			})

			It("should create Familiar's XDG config directory if it does not exist", func() {
			})

			It("should return an error if there is an issue creating Familiar's XDG config directory", func() {
			})

			It("should return an error if there is an issue creating the config location file", func() {
			})
		})
	})

	Describe("GetConfig", func() {
		It("should return the contents of the config file as a pointer to a Config struct", func() {
		})

		It("should return a new Config struct if the config file does not yet exist", func() {
		})

		It("should create the config file using a new Config struct if it does not exist", func() {
		})

		It("should close the config file after writing to it", func() {
		})

		It("should return an error if the config file location is not set", func() {
		})

		It("should return an error if there is an issue converting the file contents from YAML to a Config struct",
			func() {
			})

		It("should return an error if the config file exists but there is an issue reading it", func() {
		})

		It("should return an error if it needs to create the config file but is unable to", func() {
		})
	})

	Describe("SetConfig", func() {
		It("should write the given configuration to the config file as YAML", func() {
		})

		It("should close the config file after writing to it", func() {
		})

		It("should return an error if the config file location is not set", func() {
		})

		It("should return an error if there is an issue creating or updating the config file", func() {
		})

		It("should return an error if there is an issue converting the Config struct to YAML", func() {
		})
	})
})
