package packagemanagers

import "fmt"

var scoopPackageManager = &ScoopPackageManager{}

var packageManagers = map[string]PackageManager{
	scoopPackageManager.Name(): scoopPackageManager,
}

// GetPackageManager returns the package manager with the given name, if it exists.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager.
func GetPackageManager(packageManagerName string) (PackageManager, error) {
	packageManager, isPresent := packageManagers[packageManagerName]
	if !isPresent {
		return nil, fmt.Errorf("package manager not valid")
	}

	return packageManager, nil
}
