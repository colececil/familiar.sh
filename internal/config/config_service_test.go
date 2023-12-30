package config_test

import (
	"bytes"
	"fmt"
	. "github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"strings"
)

var _ = Describe("ConfigService", func() {
	const configHomeLocation = "/home/user/.config"
	const appDirectoryName = "io.colececil.familiar"
	const configLocationFileName = "config_location"
	const configLocation = "/path/to/config.yml"

	var fileSystemServiceDouble *test.FileSystemServiceDouble
	var outputWriterDouble *bytes.Buffer
	var configService *ConfigService

	BeforeEach(func() {
		fileSystemServiceDouble = test.NewFileSystemServiceDouble()
		fileSystemServiceDouble.SetXdgConfigHome(configHomeLocation)
		outputWriterDouble = new(bytes.Buffer)
		configService = NewConfigService(fileSystemServiceDouble, outputWriterDouble)
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
		var expectedFileContent string

		BeforeEach(func() {
			configFileDir := fileSystemServiceDouble.Dir(configLocation)
			file, _ := fileSystemServiceDouble.CreateFile(configFileDir)
			_ = file.Close()

			expectedFileContent, _ = fileSystemServiceDouble.Abs(configLocation)
			expectedFileContent += "\n"
		})

		When("the config location file exists", func() {
			BeforeEach(func() {
				file, _ := fileSystemServiceDouble.CreateFile(configHomeLocation + "/" + appDirectoryName + "/" +
					configLocationFileName)
				_ = file.Close()
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

			It("should return an error if there is an issue determining the absolute representation of the given path",
				func() {
					fileSystemServiceDouble.ReturnErrorFromMethod("Abs", configLocation)
					err := configService.SetConfigLocation(configLocation)
					Expect(err.Error()).To(Equal("unable to parse the given path"))
				})

			It("should return an error if there is an issue checking if the directory exists", func() {
				dir := fileSystemServiceDouble.Dir(configLocation)
				fileSystemServiceDouble.ReturnErrorFromMethod("FileExists", dir)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal(fmt.Sprintf("error checking existence of directory \"%s\"", dir)))
			})

			It("should return an error if there is an issue updating the config location file", func() {
				fileSystemServiceDouble.ReturnErrorFromMethod("CreateFile", configHomeLocation+"/"+appDirectoryName+"/"+
					configLocationFileName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating file"))
			})
		})

		When("the config location file does not exist", func() {
			It("should create the config location file", func() {
				err := configService.SetConfigLocation(configLocation)
				fileContentBytes, _ := fileSystemServiceDouble.ReadFile(configHomeLocation + "/" +
					appDirectoryName + "/" + configLocationFileName)
				fileContent := string(fileContentBytes)

				Expect(err).To(BeNil())
				Expect(fileContent).To(Equal(expectedFileContent))
			})

			It("should create Familiar's XDG config directory if it does not exist", func() {
				err := configService.SetConfigLocation(configLocation)
				Expect(err).To(BeNil())

				_, err = fileSystemServiceDouble.ReadFile(configHomeLocation + "/" + appDirectoryName)
				Expect(err).To(BeNil())
			})

			It("should return an error if there is an issue creating Familiar's XDG config directory", func() {
				fileSystemServiceDouble.ReturnErrorFromMethod("CreateDirectory",
					configHomeLocation+"/"+appDirectoryName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating directory"))
			})

			It("should return an error if there is an issue creating the config location file", func() {
				fileSystemServiceDouble.ReturnErrorFromMethod("CreateFile", configHomeLocation+"/"+appDirectoryName+"/"+
					configLocationFileName)
				err := configService.SetConfigLocation(configLocation)
				Expect(err.Error()).To(Equal("error creating file"))
			})
		})
	})

	Describe("GetConfig", func() {
		const packageManager1Name = "packageManager1"
		const packageManager2Name = "packageManager2"
		const packageManager3Name = "packageManager3"
		const package1Name = "package1"
		const package2Name = "package2"
		const package3Name = "package3"

		var configLocationFile io.WriteCloser
		var configFile io.WriteCloser
		var packageManagerRegistry packagemanagers.PackageManagerRegistry

		BeforeEach(func() {
			configLocationFile, _ = fileSystemServiceDouble.CreateFile(configHomeLocation + "/" + appDirectoryName +
				"/" + configLocationFileName)
			_, _ = configLocationFile.Write([]byte(configLocation))
			_ = configLocationFile.Close()

			configFile, _ = fileSystemServiceDouble.CreateFile(configLocation)

			packageManagerRegistry = packagemanagers.NewPackageManagerRegistry([]packagemanagers.PackageManager{
				test.NewPackageManagerDouble(packageManager1Name),
				test.NewPackageManagerDouble(packageManager2Name),
				test.NewPackageManagerDouble(packageManager3Name),
			})
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
packageManagers: []`,
				func(expectedConfig *Config) {},
			),
			Entry("when the config contains package managers with no packages",
				`version: 1
files: []
scripts: []
packageManagers:
    - name: packageManager1
      packages: []
    - name: packageManager2
      packages: []`,
				func(expectedConfig *Config) {
					_ = expectedConfig.AddPackageManager(packageManager1Name, packageManagerRegistry)
					_ = expectedConfig.AddPackageManager(packageManager2Name, packageManagerRegistry)
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
            version: 1.0.1`,
				func(expectedConfig *Config) {
					_ = expectedConfig.AddPackageManager(packageManager1Name, packageManagerRegistry)
					_ = expectedConfig.AddPackage(packageManager1Name, package1Name,
						packagemanagers.NewVersion("1.0.0"))
					_ = expectedConfig.AddPackage(packageManager1Name, package2Name,
						packagemanagers.NewVersion("1.1.1"))

					_ = expectedConfig.AddPackageManager(packageManager2Name, packageManagerRegistry)

					_ = expectedConfig.AddPackageManager(packageManager3Name, packageManagerRegistry)
					_ = expectedConfig.AddPackage(packageManager3Name, package3Name,
						packagemanagers.NewVersion("1.0.1"))
				},
			),
		)

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
