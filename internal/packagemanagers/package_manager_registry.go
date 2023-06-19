package packagemanagers

import (
	"fmt"
)

type PackageManagerRegistry map[string]PackageManager

// NewPackageManagerRegistry returns a new instance of PackageManagerRegistry.
func NewPackageManagerRegistry(scoopPackageManager *ScoopPackageManager) PackageManagerRegistry {
	return PackageManagerRegistry{
		scoopPackageManager.Name(): scoopPackageManager,
	}
}

// GetAllPackageManagers returns a slice containing all package managers.
func (packageManagerRegistry PackageManagerRegistry) GetAllPackageManagers() []PackageManager {
	var packageManagersSlice []PackageManager
	for _, packageManager := range packageManagerRegistry {
		packageManagersSlice = append(packageManagersSlice, packageManager)
	}
	return packageManagersSlice
}

// GetPackageManager returns the package manager with the given name, if it exists.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
func (packageManagerRegistry PackageManagerRegistry) GetPackageManager(packageManagerName string) (PackageManager,
	error) {
	packageManager, isPresent := packageManagerRegistry[packageManagerName]
	if !isPresent {
		return nil, fmt.Errorf("package manager not valid")
	}

	return packageManager, nil
}
