package packagemanagers

import "errors"

// ChocolateyPackageManager implements the PackageManager interface for the Chocolatey package manager.
type ChocolateyPackageManager struct{}

// NewChocolateyPackageManager returns a new instance of ChocolateyPackageManager.
func NewChocolateyPackageManager() *ChocolateyPackageManager {
	return &ChocolateyPackageManager{}
}

// Name implements PackageManager.Name by returning "chocolatey".
func (c *ChocolateyPackageManager) Name() string {
	return "chocolatey"
}

// Order implements PackageManager.Order by returning 2.
func (c *ChocolateyPackageManager) Order() int {
	return 2
}

// IsSupported implements PackageManager.IsSupported by returning false.
func (c *ChocolateyPackageManager) IsSupported() bool {
	// TODO: Implement this function.
	return false
}

// IsInstalled implements PackageManager.IsInstalled by returning an error.
func (c *ChocolateyPackageManager) IsInstalled() (bool, error) {
	// TODO: Implement this function.
	return false, errors.New("not implemented")
}

// Install implements PackageManager.Install by returning an error.
func (c *ChocolateyPackageManager) Install() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// Update implements PackageManager.Update by returning an error.
func (c *ChocolateyPackageManager) Update() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// Uninstall implements PackageManager.Uninstall by returning an error.
func (c *ChocolateyPackageManager) Uninstall() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// InstalledPackages implements PackageManager.InstalledPackages by returning an error.
func (c *ChocolateyPackageManager) InstalledPackages() ([]*Package, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// InstallPackage implements PackageManager.InstallPackage by returning an error.
func (c *ChocolateyPackageManager) InstallPackage(packageName string, version *Version) (*Version, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// UpdatePackage implements PackageManager.UpdatePackage by returning an error.
func (c *ChocolateyPackageManager) UpdatePackage(packageName string, version *Version) (*Version, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// UninstallPackage implements PackageManager.UninstallPackage by returning an error.
func (c *ChocolateyPackageManager) UninstallPackage(packageName string) error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}
