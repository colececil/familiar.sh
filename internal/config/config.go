package config

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"gopkg.in/yaml.v3"
	"strings"
)

// Config represents the contents of the config file.
type Config struct {
	Version         int                        `yaml:"version"`
	Files           []ConfiguredFile           `yaml:"files"`
	Scripts         []ConfiguredScript         `yaml:"scripts"`
	PackageManagers []ConfiguredPackageManager `yaml:"packageManagers"`
}

// ConfiguredFile represents a file managed by Familiar.sh.
type ConfiguredFile struct {
	SourcePath       string                      `yaml:"sourcePath"`
	DestinationPath  string                      `yaml:"destinationPath,omitempty"`
	OperatingSystems []ConfiguredOperatingSystem `yaml:"operatingSystems,omitempty"`
}

// ConfiguredScript represents a script managed by Familiar.sh.
type ConfiguredScript struct {
	SourcePath       string                      `yaml:"sourcePath"`
	OperatingSystems []ConfiguredOperatingSystem `yaml:"operatingSystems,omitempty"`
}

// ConfiguredPackageManager represents a package manager installed by Familiar.sh.
type ConfiguredPackageManager struct {
	Name     string              `yaml:"name"`
	Packages []ConfiguredPackage `yaml:"packages"`
}

// ConfiguredPackage represents a package installed by a specific package manager.
type ConfiguredPackage struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// ConfiguredOperatingSystem represents an OS that a ConfiguredFile or ConfiguredScript is used in.
type ConfiguredOperatingSystem struct {
	Name            string `yaml:"name"`
	DestinationPath string `yaml:"destinationPath,omitempty"`
}

// YamlString returns the Config contents as a YAML string.
func (config *Config) YamlString() (string, error) {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

// NewConfig creates a new instance of Config.
func NewConfig() *Config {
	return &Config{
		Version:         1,
		Files:           []ConfiguredFile{},
		Scripts:         []ConfiguredScript{},
		PackageManagers: []ConfiguredPackageManager{},
	}
}

// AddPackageManager adds the given package manager to the Config.
//
// It throws an error under the following conditions:
//   - The given package manager is not a valid package manager.
//   - The given package manager is already in the Config.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to add.
//   - packageManagerRegistry: The package manager registry to use for validating the package manager name.
func (config *Config) AddPackageManager(packageManagerName string, packageManagerRegistry packagemanagers.PackageManagerRegistry) error {
	if _, err := packageManagerRegistry.GetPackageManager(packageManagerName); err != nil {
		return fmt.Errorf("package manager not valid")
	}

	for i := range config.PackageManagers {
		if config.PackageManagers[i].Name == packageManagerName {
			return fmt.Errorf("package manager already present")
		}
	}

	packageManager := ConfiguredPackageManager{
		Name:     packageManagerName,
		Packages: []ConfiguredPackage{},
	}
	config.PackageManagers = append(config.PackageManagers, packageManager)
	return nil
}

// RemovePackageManager removes the given package manager from the Config. If the given package manager is not present
// in the Config, it throws an error.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to remove.
func (config *Config) RemovePackageManager(packageManagerName string) error {
	var filteredPackageManagers []ConfiguredPackageManager
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
}

// AddPackage updates the Config to add the given version of the given package under the given package manager.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the Config.
//   - The given package is already in the Config under the given package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to add.
//   - packageVersion: The version of the package to add.
func (config *Config) AddPackage(packageManagerName string, packageName string,
	packageVersion *packagemanagers.Version) error {
	var matchingPackageManager *ConfiguredPackageManager
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

	newPackage := ConfiguredPackage{Name: packageName, Version: packageVersion.VersionString}
	matchingPackageManager.Packages = append(matchingPackageManager.Packages, newPackage)
	return nil
}

// UpdatePackage updates the Config to change the given package under the given package manager to the given version.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the Config.
//   - The given package is not in the Config under the given package manager.
//   - The given version is the same as the version of the package in the Config.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to update.
//   - packageVersion: The version of the package to update.
func (config *Config) UpdatePackage(packageManagerName string, packageName string,
	packageVersion *packagemanagers.Version) error {
	var matchingPackageManager *ConfiguredPackageManager
	for i := range config.PackageManagers {
		if config.PackageManagers[i].Name == packageManagerName {
			matchingPackageManager = &config.PackageManagers[i]
			break
		}
	}

	if matchingPackageManager == nil {
		return fmt.Errorf("package manager not present")
	}

	var matchingPackage *ConfiguredPackage
	for i := range matchingPackageManager.Packages {
		if matchingPackageManager.Packages[i].Name == packageName {
			matchingPackage = &matchingPackageManager.Packages[i]
			break
		}
	}

	if matchingPackage == nil {
		return fmt.Errorf("package not present")
	}

	if matchingPackage.Version == packageVersion.VersionString {
		return fmt.Errorf("package already set to given version")
	}

	matchingPackage.Version = packageVersion.VersionString
	return nil
}

// RemovePackage updates the Config to remove the given package under the given package manager.
//
// It throws an error under the following conditions:
//   - The given package manager is not in the Config.
//   - The given package is not in the Config under the given package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
//   - packageName: The name of the package to add.
func (config *Config) RemovePackage(packageManagerName string, packageName string) error {
	var matchingPackageManager *ConfiguredPackageManager
	for i := range config.PackageManagers {
		if config.PackageManagers[i].Name == packageManagerName {
			matchingPackageManager = &config.PackageManagers[i]
			break
		}
	}

	if matchingPackageManager == nil {
		return fmt.Errorf("package manager not present")
	}

	var filteredPackages []ConfiguredPackage
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
}
