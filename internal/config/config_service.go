package config

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

const appDirectoryName = "io.colececil.familiar"
const configLocationFileName = "config_location"
const configLocationNotSetError = "The location of Familiar's shared config file has not yet been set. Please set it " +
	"using \"familiar config location <path>\", for more details, execute \"familiar help config\"."

// ConfigService is a service that manages the shared configuration file.
type ConfigService struct {
	fileSystemService system.FileSystemService
	outputWriter      io.Writer
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(fileSystemService system.FileSystemService, outputWriter io.Writer) *ConfigService {
	return &ConfigService{
		fileSystemService: fileSystemService,
		outputWriter:      outputWriter,
	}
}

// GetConfigLocation returns the location of the shared configuration file, as stored in the "config_location" file in
// Familiar's XDG config directory. If the "config_location" file does not exist or is empty, an error is returned.
func (service *ConfigService) GetConfigLocation() (string, error) {
	configDir := service.fileSystemService.GetXdgConfigHome()

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName
	bytes, err := service.fileSystemService.ReadFile(configLocationFilePath)
	if err != nil {
		return "", fmt.Errorf(configLocationNotSetError)
	}

	path := strings.TrimSpace(string(bytes))
	if len(path) == 0 {
		return "", fmt.Errorf(configLocationNotSetError)
	}

	return path, nil
}

// SetConfigLocation writes the given path to the "config_location" file in Familiar's XDG config directory. If the
// Familiar XDG config directory or the "config_location" file do not exist, they will be created. If the
// "config_location" file already exists, it will be overwritten.
//
// If the directory of the file specified by the given path does not exist, or if the file specified by the path does
// not have a YAML file extension, an error is returned.
func (service *ConfigService) SetConfigLocation(path string) error {
	absolutePath, err := service.fileSystemService.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to parse the given path")
	}

	ext := service.fileSystemService.Ext(absolutePath)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("invalid file extension \"%s\": expected \".yml\" or \".yaml\"", ext)
	}

	dirPath := service.fileSystemService.Dir(absolutePath)
	directoryExists, err := service.fileSystemService.FileExists(dirPath)
	if err != nil {
		return fmt.Errorf("error checking existence of directory \"%s\"", dirPath)
	}
	if !directoryExists {
		return fmt.Errorf("directory \"%s\" does not exist", dirPath)
	}

	configDir := service.fileSystemService.GetXdgConfigHome()

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName

	err = service.fileSystemService.CreateDirectory(configDir+"/"+appDirectoryName, 0660)
	if err != nil {
		return err
	}

	file, err := service.fileSystemService.CreateFile(configLocationFilePath)
	if err != nil {
		return err
	}
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	_, err = fmt.Fprintln(file, absolutePath)
	if err == nil {
		_, err = fmt.Fprintf(service.outputWriter, "The config file location has been set to \"%s\".\n", absolutePath)
	}
	return err
}

// GetConfig returns the contents of the config file as a pointer to a Config struct. If the config file does not yet
// exist, a new Config struct is created and written to the file before being returned.
//
// An error is returned if the config file location is not set, or if there is an issue unmarshalling the config file
// to a Config struct.
func (service *ConfigService) GetConfig() (*Config, error) {
	configLocation, err := service.GetConfigLocation()
	if err != nil {
		return nil, err
	}

	bytes, err := service.fileSystemService.ReadFile(configLocation)
	if err != nil {
		if os.IsNotExist(err) {
			newConfig := NewConfig()
			if err = service.SetConfig(newConfig); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var config Config
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetConfig writes the given configuration to the config file as YAML. An error is returned if the config file location
// is not set, or if there is an issue creating or updating the config file.
func (service *ConfigService) SetConfig(config *Config) error {
	configLocation, err := service.GetConfigLocation()
	if err != nil {
		return err
	}

	file, err := service.fileSystemService.CreateFile(configLocation)
	if err != nil {
		return err
	}
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}
