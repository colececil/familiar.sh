package packagemanagers

import (
	"fmt"
	"slices"
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
func (r PackageManagerRegistry) GetAllPackageManagers() []PackageManager {
	var packageManagersSlice []PackageManager
	for _, packageManager := range r {
		packageManagersSlice = append(packageManagersSlice, packageManager)
	}
	slices.SortFunc(packageManagersSlice, func(packageManager1, packageManager2 PackageManager) int {
		return packageManager1.Order() - packageManager2.Order()
	})
	return packageManagersSlice
}

// GetPackageManager returns the package manager with the given name, if it exists.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
func (r PackageManagerRegistry) GetPackageManager(packageManagerName string) (PackageManager,
	error) {
	packageManager, isPresent := r[packageManagerName]
	if !isPresent {
		return nil, fmt.Errorf("package manager not valid")
	}

	return packageManager, nil
}

func (r PackageManagerRegistry) validate() {
	panicMessage := "package manager registry does not contain the expected package managers"
	expectedPackageManagers := []string{"scoop", "chocolatey", "homebrew"}

	if len(r) != len(expectedPackageManagers) {
		panic(panicMessage)
	}

	for _, expectedPackageManager := range expectedPackageManagers {
		if r[expectedPackageManager] == nil {
			panic(panicMessage)
		}
	}

	orderNumbersEncountered := make(map[int]bool)
	for _, packageManager := range r {
		order := packageManager.Order()
		if order < 1 || order > len(r) || orderNumbersEncountered[order] {
			panic(panicMessage)
		}
		orderNumbersEncountered[order] = true
	}
}
