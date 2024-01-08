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
	fileSystem   FileSystem
	outputWriter io.Writer
}

// NewConfigService creates a new instance of ConfigService.
func NewConfigService(fileSystem FileSystem, outputWriter io.Writer) *ConfigService {

	return &ConfigService{
		fileSystem:   fileSystem,
		outputWriter: outputWriter,
	}
}

// FileSystem is an interface that provides necessary methods to interact with the underlying file system.
type FileSystem interface {
	system.XdgConfigHomeGetter
	system.AbsPathConverter
	system.PathDirGetter
	system.FileExtensionGetter
	system.DirCreator
	system.FileExistenceChecker
	system.FileReader
	system.FileCreator
}

// GetConfigLocation returns the location of the shared configuration file, as stored in the "config_location" file in
// Familiar's XDG config directory. If the "config_location" file does not exist or is empty, an error is returned.
func (s *ConfigService) GetConfigLocation() (string, error) {
	configDir := s.fileSystem.GetXdgConfigHome()

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName
	bytes, err := s.fileSystem.ReadFile(configLocationFilePath)
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
func (s *ConfigService) SetConfigLocation(path string) error {
	absolutePath, err := s.fileSystem.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to parse the given path")
	}

	ext := s.fileSystem.Ext(absolutePath)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("invalid file extension \"%s\": expected \".yml\" or \".yaml\"", ext)
	}

	dirPath := s.fileSystem.Dir(absolutePath)
	directoryExists, err := s.fileSystem.FileExists(dirPath)
	if err != nil {
		return fmt.Errorf("error checking existence of directory \"%s\"", dirPath)
	}
	if !directoryExists {
		return fmt.Errorf("directory \"%s\" does not exist", dirPath)
	}

	configDir := s.fileSystem.GetXdgConfigHome()

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName

	err = s.fileSystem.CreateDirectory(configDir+"/"+appDirectoryName, 0660)
	if err != nil {
		return err
	}

	file, err := s.fileSystem.CreateFile(configLocationFilePath)
	if err != nil {
		return err
	}
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	_, err = fmt.Fprintln(file, absolutePath)
	if err == nil {
		_, err = fmt.Fprintf(s.outputWriter, "The config file location has been set to \"%s\".\n", absolutePath)
	}
	return err
}

// GetConfig returns the contents of the config file as a pointer to a Config struct. If the config file does not yet
// exist, a new Config struct is created and written to the file before being returned.
//
// An error is returned if the config file location is not set, or if there is an issue unmarshalling the config file
// to a Config struct.
func (s *ConfigService) GetConfig() (*Config, error) {
	configLocation, err := s.GetConfigLocation()
	if err != nil {
		return nil, err
	}

	bytes, err := s.fileSystem.ReadFile(configLocation)
	if err != nil {
		if os.IsNotExist(err) {
			newConfig := NewConfig()
			if err = s.SetConfig(newConfig); err != nil {
				return nil, err
			}
			return newConfig, nil
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
func (s *ConfigService) SetConfig(config *Config) error {
	configLocation, err := s.GetConfigLocation()
	if err != nil {
		return err
	}

	file, err := s.fileSystem.CreateFile(configLocation)
	if err != nil {
		return err
	}
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}
