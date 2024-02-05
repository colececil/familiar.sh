package config_test

import (
	"bytes"
	"fmt"
	. "github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
	"io"
	"strings"
)

var _ = Describe("ConfigService", func() {
	const configHomeLocation = "/home/user/.config"
	const appDirectoryName = "io.colececil.familiar"
	const configLocationFileName = "config_location"
	const configLocation = "/path/to/config.yml"

	var fileSystemDouble *test.FileSystemDouble
	var outputWriterDouble *bytes.Buffer
	var configService *ConfigService

	BeforeEach(func() {
		fileSystemDouble = test.NewFileSystemDouble()
		fileSystemDouble.SetXdgConfigHome(configHomeLocation)
		outputWriterDouble = new(bytes.Buffer)
		configService = NewConfigService(fileSystemDouble, outputWriterDouble)
	})

	Describe("GetConfigLocation", func() {
		It("should return the config location defined in the config location file", func() {
			file, _ := fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
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
			file, _ := fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
				configLocationFileName)
			_ = file.Close()

			location, err := configService.GetConfigLocation()

			Expect(err).To(Equal(expectedError))
			Expect(location).To(Equal(""))
		})
	})

	Describe("SetConfigLocation", func() {
		var expectedFileContent string

		BeforeEach(func() {
			configFileDir := fileSystemDouble.Dir(configLocation)
			file, _ := fileSystemDouble.CreateFile(configFileDir)
			_ = file.Close()

			expectedFileContent, _ = fileSystemDouble.Abs(configLocation)
			expectedFileContent += "\n"
		})

		When("the config location file exists", func() {
			BeforeEach(func() {
				file, _ := fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
					configLocationFileName)
				_ = file.Close()
			})

			It("should write the given path to the config location file", func() {
				err := configService.SetConfigLocation(configLocation)
				fileContentBytes, _ := fileSystemDouble.ReadFile(configHomeLocation + "/" +
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
				file, _ := fileSystemDouble.GetCreatedFile(configHomeLocation + "/" + appDirectoryName + "/" +
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

			It("should return an error if there is an issue determining the absolute representation of the given path",
				func() {
					fileSystemDouble.ReturnErrorFromMethod("Abs", configLocation)
					err := configService.SetConfigLocation(configLocation)
					Expect(err.Error()).To(Equal("unable to parse the given path"))
				})

			It("should return an error if there is an issue checking if the directory exists", func() {
				dir := fileSystemDouble.Dir(configLocation)
				fileSystemDouble.ReturnErrorFromMethod("FileExists", dir)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal(fmt.Sprintf("error checking existence of directory \"%s\"", dir)))
			})

			It("should return an error if there is an issue updating the config location file", func() {
				fileSystemDouble.ReturnErrorFromMethod("CreateFile", configHomeLocation+"/"+appDirectoryName+"/"+
					configLocationFileName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating file"))
			})
		})

		When("the config location file does not exist", func() {
			It("should create the config location file", func() {
				err := configService.SetConfigLocation(configLocation)
				fileContentBytes, _ := fileSystemDouble.ReadFile(configHomeLocation + "/" +
					appDirectoryName + "/" + configLocationFileName)
				fileContent := string(fileContentBytes)

				Expect(err).To(BeNil())
				Expect(fileContent).To(Equal(expectedFileContent))
			})

			It("should create Familiar's XDG config directory if it does not exist", func() {
				err := configService.SetConfigLocation(configLocation)
				Expect(err).To(BeNil())

				_, err = fileSystemDouble.ReadFile(configHomeLocation + "/" + appDirectoryName)
				Expect(err).To(BeNil())
			})

			It("should return an error if there is an issue creating Familiar's XDG config directory", func() {
				fileSystemDouble.ReturnErrorFromMethod("CreateDirectory",
					configHomeLocation+"/"+appDirectoryName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating directory"))
			})

			It("should return an error if there is an issue creating the config location file", func() {
				fileSystemDouble.ReturnErrorFromMethod("CreateFile", configHomeLocation+"/"+appDirectoryName+"/"+
					configLocationFileName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating file"))
			})
		})
	})

	Describe("GetConfig", func() {
		const scoop = "scoop"
		const chocolatey = "chocolatey"
		const homebrew = "homebrew"
		const package1Name = "package1"
		const package2Name = "package2"
		const package3Name = "package3"

		var configLocationFile io.WriteCloser
		var configFile io.WriteCloser
		var packageManagerRegistry packagemanagers.PackageManagerRegistry

		BeforeEach(func() {
			mock.SetUp(GinkgoT())

			scoopPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(scoopPackageManager.Name()).ThenReturn(scoop)
			mock.WhenSingle(scoopPackageManager.Order()).ThenReturn(1)

			chocolateyPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(chocolateyPackageManager.Name()).ThenReturn(chocolatey)
			mock.WhenSingle(chocolateyPackageManager.Order()).ThenReturn(2)

			homebrewPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(homebrewPackageManager.Name()).ThenReturn(homebrew)
			mock.WhenSingle(homebrewPackageManager.Order()).ThenReturn(3)

			packageManagerRegistry = packagemanagers.NewPackageManagerRegistry([]packagemanagers.PackageManager{
				scoopPackageManager,
				chocolateyPackageManager,
				homebrewPackageManager,
			})

			configLocationFile, _ = fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
				"/" + configLocationFileName)
			_, _ = configLocationFile.Write([]byte(configLocation))
			_ = configLocationFile.Close()

			configFile, _ = fileSystemDouble.CreateFile(configLocation)
		})

		DescribeTable("should return the contents of the config file as a pointer to a Config struct",
			func(fileContent string, expectedConfigSetupFunc func(expectedConfig *Config)) {
				_, _ = configFile.Write([]byte(fileContent))
				_ = configFile.Close()

				expectedConfig := NewConfig()
				expectedConfigSetupFunc(expectedConfig)

				result, err := configService.GetConfig()

				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedConfig))
			},
			Entry("when the config is empty",
				`version: 1
files: []
scripts: []
packageManagers: []
`,
				func(expectedConfig *Config) {},
			),
			Entry("when the config contains package managers with no packages",
				`version: 1
files: []
scripts: []
packageManagers:
  - name: chocolatey
    packages: []
  - name: scoop
    packages: []
`,
				func(expectedConfig *Config) {
					_ = expectedConfig.AddPackageManager(scoop, packageManagerRegistry)
					_ = expectedConfig.AddPackageManager(chocolatey, packageManagerRegistry)
				},
			),
			Entry("when the config contains package managers with packages",
				`version: 1
files: []
scripts: []
packageManagers:
  - name: chocolatey
    packages: []
  - name: homebrew
    packages:
      - name: package3
        version: 1.0.1
  - name: scoop
    packages:
      - name: package1
        version: 1.0.0
      - name: package2
        version: 1.1.1
`,
				func(expectedConfig *Config) {
					_ = expectedConfig.AddPackageManager(scoop, packageManagerRegistry)
					_ = expectedConfig.AddPackage(scoop, package1Name,
						packagemanagers.NewVersion("1.0.0"))
					_ = expectedConfig.AddPackage(scoop, package2Name,
						packagemanagers.NewVersion("1.1.1"))

					_ = expectedConfig.AddPackageManager(chocolatey, packageManagerRegistry)

					_ = expectedConfig.AddPackageManager(homebrew, packageManagerRegistry)
					_ = expectedConfig.AddPackage(homebrew, package3Name,
						packagemanagers.NewVersion("1.0.1"))
				},
			),
		)

		When("the config file does not exist", func() {
			const otherConfigLocation = "/other/path/to/config.yml"

			BeforeEach(func() {
				configLocationFile, _ = fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
					"/" + configLocationFileName)
				_, _ = configLocationFile.Write([]byte(otherConfigLocation))
				_ = configLocationFile.Close()
			})

			It("should return a new Config struct", func() {
				expectedConfig := NewConfig()

				result, err := configService.GetConfig()

				Expect(err).To(BeNil())
				Expect(result).To(Equal(expectedConfig))
			})

			It("should create the config file using a new Config struct", func() {
				expectedFileContent := `version: 1
files: []
scripts: []
packageManagers: []
`

				_, err := configService.GetConfig()
				fileContentBytes, _ := fileSystemDouble.ReadFile(otherConfigLocation)
				fileContent := string(fileContentBytes)

				Expect(err).To(BeNil())
				Expect(fileContent).To(Equal(expectedFileContent))
			})

			It("should close the config file after writing to it", func() {
				_, err := configService.GetConfig()
				file, _ := fileSystemDouble.GetCreatedFile(otherConfigLocation)

				Expect(err).To(BeNil())
				Expect(file.IsClosed()).To(BeTrue())
			})

			It("should return an error if there is an issue creating the config file", func() {
				fileSystemDouble.ReturnErrorFromMethod("CreateFile", otherConfigLocation)
				_, err := configService.GetConfig()
				Expect(err.Error()).To(Equal("error creating file"))
			})
		})

		It("should return an error if the config file location is not set", func() {
			configLocationFile, _ = fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
				"/" + configLocationFileName)
			_, _ = configLocationFile.Write([]byte{})
			_ = configLocationFile.Close()

			_, err := configService.GetConfig()

			Expect(err.Error()).To(Equal("The location of Familiar's shared config file has not yet been set. Please " +
				"set it using \"familiar config location <path>\", for more details, execute \"familiar help " +
				"config\"."))
		})

		It("should return an error if the config file exists but there is an issue reading it", func() {
			fileSystemDouble.ReturnErrorFromMethod("ReadFile", configLocation)
			_, err := configService.GetConfig()
			Expect(err.Error()).To(Equal("error reading file"))
		})
	})

	Describe("SetConfig", func() {
		const scoop = "scoop"
		const chocolatey = "chocolatey"
		const homebrew = "homebrew"
		const package1Name = "package1"
		const package2Name = "package2"
		const package3Name = "package3"

		var configLocationFile io.WriteCloser
		var config *Config
		var packageManagerRegistry packagemanagers.PackageManagerRegistry

		BeforeEach(func() {
			mock.SetUp(GinkgoT())

			scoopPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(scoopPackageManager.Name()).ThenReturn(scoop)
			mock.WhenSingle(scoopPackageManager.Order()).ThenReturn(1)

			chocolateyPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(chocolateyPackageManager.Name()).ThenReturn(chocolatey)
			mock.WhenSingle(chocolateyPackageManager.Order()).ThenReturn(2)

			homebrewPackageManager := mock.Mock[packagemanagers.PackageManager]()
			mock.WhenSingle(homebrewPackageManager.Name()).ThenReturn(homebrew)
			mock.WhenSingle(homebrewPackageManager.Order()).ThenReturn(3)

			packageManagerRegistry = packagemanagers.NewPackageManagerRegistry([]packagemanagers.PackageManager{
				scoopPackageManager,
				chocolateyPackageManager,
				homebrewPackageManager,
			})

			configLocationFile, _ = fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
				"/" + configLocationFileName)
			_, _ = configLocationFile.Write([]byte(configLocation))
			_ = configLocationFile.Close()

			config = NewConfig()
		})

		DescribeTable("should write the given configuration to the config file as YAML",
			func(expectedFileContent string, configSetupFunc func()) {
				configSetupFunc()

				err := configService.SetConfig(config)
				fileContentBytes, _ := fileSystemDouble.ReadFile(configLocation)
				fileContent := string(fileContentBytes)

				Expect(err).To(BeNil())

				parsedExpected, err := parseYamlString(expectedFileContent)
				Expect(err).To(BeNil())
				parsedResult, err := parseYamlString(fileContent)
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
  - name: chocolatey
    packages: []
  - name: scoop
    packages: []
`,
				func() {
					_ = config.AddPackageManager(scoop, packageManagerRegistry)
					_ = config.AddPackageManager(chocolatey, packageManagerRegistry)
				},
			),
			Entry("when the config contains package managers with packages",
				`version: 1
files: []
scripts: []
packageManagers:
  - name: chocolatey
    packages: []
  - name: homebrew
    packages:
      - name: package3
        version: 1.0.1
  - name: scoop
    packages:
      - name: package1
        version: 1.0.0
      - name: package2
        version: 1.1.1
`,
				func() {
					_ = config.AddPackageManager(scoop, packageManagerRegistry)
					_ = config.AddPackage(scoop, package1Name,
						packagemanagers.NewVersion("1.0.0"))
					_ = config.AddPackage(scoop, package2Name,
						packagemanagers.NewVersion("1.1.1"))

					_ = config.AddPackageManager(chocolatey, packageManagerRegistry)

					_ = config.AddPackageManager(homebrew, packageManagerRegistry)
					_ = config.AddPackage(homebrew, package3Name,
						packagemanagers.NewVersion("1.0.1"))
				},
			),
		)

		It("should close the config file after writing to it", func() {
			err := configService.SetConfig(config)
			file, _ := fileSystemDouble.GetCreatedFile(configLocation)

			Expect(err).To(BeNil())
			Expect(file.IsClosed()).To(BeTrue())
		})

		It("should return an error if the config file location is not set", func() {
			configLocationFile, _ = fileSystemDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
				"/" + configLocationFileName)
			_, _ = configLocationFile.Write([]byte{})
			_ = configLocationFile.Close()

			err := configService.SetConfig(config)
			Expect(err.Error()).To(Equal("The location of Familiar's shared config file has not yet been set. Please " +
				"set it using \"familiar config location <path>\", for more details, execute \"familiar help " +
				"config\"."))
		})

		It("should return an error if there is an issue creating or updating the config file", func() {
			fileSystemDouble.ReturnErrorFromMethod("CreateFile", configLocation)
			err := configService.SetConfig(config)
			Expect(err.Error()).To(Equal("error creating file"))
		})
	})
})
