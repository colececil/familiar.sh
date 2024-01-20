package packagemanagers

import (
	"fmt"
	"slices"
	"strings"
)

type PackageManagerRegistry map[string]PackageManager

// NewPackageManagerRegistry returns a new instance of PackageManagerRegistry.
func NewPackageManagerRegistry(packageManagers []PackageManager) PackageManagerRegistry {
	packageManagerRegistry := make(PackageManagerRegistry)
	for _, packageManager := range packageManagers {
		packageManagerRegistry[packageManager.Name()] = packageManager
	}
	packageManagerRegistry.validate()
	return packageManagerRegistry
}

// GetAllPackageManagers returns a slice containing all package managers.
func (packageManagerRegistry PackageManagerRegistry) GetAllPackageManagers() []PackageManager {
	var packageManagersSlice []PackageManager
	for _, packageManager := range packageManagerRegistry {
		packageManagersSlice = append(packageManagersSlice, packageManager)
	}
	slices.SortFunc(packageManagersSlice, func(packageManager1, packageManager2 PackageManager) int {
		return strings.Compare(
			strings.ToLower(packageManager1.Name()),
			strings.ToLower(packageManager2.Name()),
		)
	})
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

// Todo: Implement and test this method.
func (packageManagerRegistry PackageManagerRegistry) validate() {
	panic("not implemented")
}
