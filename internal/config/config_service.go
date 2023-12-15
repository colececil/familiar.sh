package config

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const appDirectoryName = "io.colececil.familiar"
const configLocationFileName = "config_location"
const configLocationNotSetError = "The location of Familiar's shared config file has not yet been set. Please set it " +
	"using \"familiar config location <path>\", for more details, execute \"familiar help config\"."

// ConfigService is a service that manages the shared configuration file.
type ConfigService struct {
	fileSystemService system.FileSystemService
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(fileSystemService system.FileSystemService) *ConfigService {
	return &ConfigService{
		fileSystemService: fileSystemService,
	}
}

// GetConfigLocation returns the location of the shared configuration file, as stored in the "config_location" file in
// the XDG config directory. If the "config_location" file does not exist or is empty, an empty string is returned.
func (service *ConfigService) GetConfigLocation() (string, error) {
	configDir := service.fileSystemService.GetXdgConfigHome()

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName
	bytes, err := service.fileSystemService.ReadFile(configLocationFilePath)
	if err != nil {
		return "", fmt.Errorf(configLocationNotSetError)
	}

	path := strings.TrimSpace(string(bytes))
	return path, nil
}

// SetConfigLocation writes the given path to the "config_location" file in the XDG config directory. If the
// directory of the file specified by the given path does not exist, or if the file is not a YAML file, an error is
// returned.
func (service *ConfigService) SetConfigLocation(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to parse the given path")
	}

	ext := filepath.Ext(absolutePath)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("invalid file extension '%s': expected '.yml' or '.yaml'", ext)
	}

	dirPath := filepath.Dir(absolutePath)
	directoryExists, err := service.fileSystemService.FileExists(dirPath)
	if err != nil {
		return fmt.Errorf("error checking directory '%s': %w", dirPath, err)
	}
	if !directoryExists {
		return fmt.Errorf("directory '%s' does not exist", dirPath)
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
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = fmt.Fprintln(file, absolutePath)
	if err == nil {
		fmt.Println("The config file location has been set to \"" + absolutePath + "\".")
	}
	return err
}

// GetConfig returns the contents of the config file as a pointer to a Config struct.
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

// SetConfig writes the given configuration to the config file as YAML.
//
// It takes the following parameters:
//   - config: The configuration to write to the file.
func (service *ConfigService) SetConfig(config *Config) error {
	configLocation, err := service.GetConfigLocation()
	if err != nil {
		return err
	}

	file, err := service.fileSystemService.CreateFile(configLocation)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}
