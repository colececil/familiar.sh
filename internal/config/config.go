package config

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

// Config represents the contents of the config file.
type Config struct {
	Version         int              `yaml:"version"`
	Files           []File           `yaml:"files"`
	Scripts         []Script         `yaml:"scripts"`
	PackageManagers []PackageManager `yaml:"packageManagers"`
}

// File represents a file managed by Familiar.sh.
type File struct {
	SourcePath       string            `yaml:"sourcePath"`
	DestinationPath  string            `yaml:"destinationPath,omitempty"`
	OperatingSystems []OperatingSystem `yaml:"operatingSystems,omitempty"`
}

// Script represents a script managed by Familiar.sh.
type Script struct {
	SourcePath       string            `yaml:"sourcePath"`
	OperatingSystems []OperatingSystem `yaml:"operatingSystems,omitempty"`
}

// PackageManager represents a package manager installed by Familiar.sh.
type PackageManager struct {
	Name     string    `yaml:"name"`
	Packages []Package `yaml:"packages"`
}

// Package represents a package installed by a specific package manager.
type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// OperatingSystem represents an OS that a File or Script is used in.
type OperatingSystem struct {
	Name            string `yaml:"name"`
	DestinationPath string `yaml:"destinationPath,omitempty"`
}

// GetConfigContents returns the contents of the config file as a YAML string.
func GetConfigContents() (string, error) {
	config, err := readConfigFile()
	if err != nil {
		return "", err
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

// AddPackageManager adds the given package manager to the config file. If the given package manager is already in the
// config file, it throws an error.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to add.
func AddPackageManager(packageManagerName string) error {
	return updateConfigFile(func(config *Config) error {
		if _, err := packagemanagers.GetPackageManager(packageManagerName); err != nil {
			return fmt.Errorf("package manager not valid")
		}

		for i := range config.PackageManagers {
			if config.PackageManagers[i].Name == packageManagerName {
				return fmt.Errorf("package manager already present")
			}
		}

		packageManager := PackageManager{
			Name:     packageManagerName,
			Packages: []Package{},
		}
		config.PackageManagers = append(config.PackageManagers, packageManager)
		return nil
	})
}

// RemovePackageManager removes the given package manager from the config file. If the given package manager is not
// present in the config file, it throws an error.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to remove.
func RemovePackageManager(packageManagerName string) error {
	return updateConfigFile(func(config *Config) error {
		var filteredPackageManagers []PackageManager
		for i := range config.PackageManagers {
			if config.PackageManagers[i].Name != packageManagerName {
				filteredPackageManagers = append(filteredPackageManagers, config.PackageManagers[i])
			}
		}

		if len(filteredPackageManagers) == len(config.PackageManagers) {
			return fmt.Errorf("package manager not present")
		}

		config.PackageManagers = filteredPackageManagers
		return nil
	})
}

// AddPackage updates the config file to add the given version of the given package under the given package manager.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the config file.
//   - The given package is already in the config file under the given package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to add.
//   - packageVersion: The version of the package to add.
func AddPackage(packageManagerName string, packageName string, packageVersion string) error {
	return updateConfigFile(func(config *Config) error {
		var matchingPackageManager *PackageManager
		for i := range config.PackageManagers {
			if config.PackageManagers[i].Name == packageManagerName {
				matchingPackageManager = &config.PackageManagers[i]
				break
			}
		}

		if matchingPackageManager == nil {
			return fmt.Errorf("package manager not present")
		}

		for i := range matchingPackageManager.Packages {
			if matchingPackageManager.Packages[i].Name == packageName {
				return fmt.Errorf("package already present")
			}
		}

		newPackage := Package{Name: packageName, Version: packageVersion}
		matchingPackageManager.Packages = append(matchingPackageManager.Packages, newPackage)
		return nil
	})
}

// UpdatePackage updates the config file to change the given package under the given package manager to the given
// version.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the config file.
//   - The given package is not in the config file under the given package manager.
//   - The given version is the same as the version of the package in the config file.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to add.
//   - packageVersion: The version of the package to add.
func UpdatePackage(packageManagerName string, packageName string, packageVersion string) error {
	return updateConfigFile(func(config *Config) error {
		var matchingPackageManager *PackageManager
		for i := range config.PackageManagers {
			if config.PackageManagers[i].Name == packageManagerName {
				matchingPackageManager = &config.PackageManagers[i]
				break
			}
		}

		if matchingPackageManager == nil {
			return fmt.Errorf("package manager not present")
		}

		var matchingPackage *Package
		for i := range matchingPackageManager.Packages {
			if matchingPackageManager.Packages[i].Name == packageName {
				matchingPackage = &matchingPackageManager.Packages[i]
				break
			}
		}

		if matchingPackage == nil {
			return fmt.Errorf("package not present")
		}

		if matchingPackage.Version == packageVersion {
			return fmt.Errorf("package already set to given version")
		}

		matchingPackage.Version = packageVersion
		return nil
	})
}

// RemovePackage updates the config file to remove the given package under the given package manager.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the config file.
//   - The given package is not in the config file under the given package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to add.
func RemovePackage(packageManagerName string, packageName string) error {
	return updateConfigFile(func(config *Config) error {
		var matchingPackageManager *PackageManager
		for i := range config.PackageManagers {
			if config.PackageManagers[i].Name == packageManagerName {
				matchingPackageManager = &config.PackageManagers[i]
				break
			}
		}

		if matchingPackageManager == nil {
			return fmt.Errorf("package manager not present")
		}

		var filteredPackages []Package
		for i := range matchingPackageManager.Packages {
			if matchingPackageManager.Packages[i].Name != packageName {
				filteredPackages = append(filteredPackages, matchingPackageManager.Packages[i])
			}
		}

		if len(filteredPackages) == len(matchingPackageManager.Packages) {
			return fmt.Errorf("package not present")
		}

		matchingPackageManager.Packages = filteredPackages
		return nil
	})
}

// readConfigFile returns the contents of the config file as a pointer to Config struct.
//
// If the config file has not yet been created, it first creates and initializes it.
func readConfigFile() (*Config, error) {
	configLocation, err := GetConfigLocation()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(configLocation)
	if err != nil {
		if os.IsNotExist(err) {
			newConfig := &Config{Version: 1, Files: []File{}, Scripts: []Script{}, PackageManagers: []PackageManager{}}
			err = writeConfigFile(newConfig)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// writeConfigFile writes the given configuration to the config file as YAML.
//
// It takes the following parameters:
//   - config: The configuration to write to the file.
func writeConfigFile(config *Config) error {
	configLocation, err := GetConfigLocation()
	if err != nil {
		return err
	}

	file, err := os.Create(configLocation)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(config)
}

// updateConfigFile updates the contents of the config file using the given function. It passes the existing contents of
// the config file into the given function, and then it saves the updated contents to the file.
//
// It takes the following parameters:
//   - updater: A function that takes in a pointer to an instance of Config and modifies it in some way.
func updateConfigFile(updater func(*Config) error) error {
	config, err := readConfigFile()
	if err != nil {
		return err
	}

	err = updater(config)
	if err != nil {
		return err
	}

	err = writeConfigFile(config)
	return err
}
