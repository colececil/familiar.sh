package config

import (
	"fmt"
	"github.com/adrg/xdg"
	"os"
	"path/filepath"
	"strings"
)

const appDirectoryName = "io.colececil.familiar"
const configLocationFileName = "config_location"
const configLocationNotSetError = "The location of Familiar's shared config file has not yet been set. Please set it " +
	"using \"familiar config location <path>\", for more details, execute \"familiar help config\"."

// GetConfigLocation returns the location of the shared configuration file, as stored in the "config_location" file in
// the XDG config directory. If the "config_location" file does not exist or is empty, an empty string is returned.
func GetConfigLocation() (string, error) {
	configDir := xdg.ConfigHome

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName
	bytes, err := os.ReadFile(configLocationFilePath)
	if err != nil {
		return "", fmt.Errorf(configLocationNotSetError)
	}

	path := strings.TrimSpace(string(bytes))
	return path, nil
}

// SetConfigLocation writes the given path to the "config_location" file in the XDG config directory. If the
// directory of the file specified by the given path does not exist, or if the file is not a YAML file, an error is
// returned.
func SetConfigLocation(path string) error {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("unable to parse the given path")
	}

	ext := filepath.Ext(absolutePath)
	if ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("invalid file extension '%s': expected '.yml' or '.yaml'", ext)
	}

	dir := filepath.Dir(absolutePath)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory '%s' does not exist", dir)
		}
		return fmt.Errorf("error checking directory '%s': %w", dir, err)
	}

	configDir := xdg.ConfigHome

	configLocationFilePath := configDir + "/" + appDirectoryName + "/" + configLocationFileName

	err = os.MkdirAll(configDir+"/"+appDirectoryName, 0700)
	if err != nil {
		return err
	}

	file, err := os.Create(configLocationFilePath)
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
