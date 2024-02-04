package packagemanagers

import "errors"

// HomebrewPackageManager implements the PackageManager interface for the Homebrew package manager.
type HomebrewPackageManager struct{}

// NewHomebrewPackageManager returns a new instance of HomebrewPackageManager.
func NewHomebrewPackageManager() *HomebrewPackageManager {
	return &HomebrewPackageManager{}
}

// Name implements PackageManager.Name by returning "chocolatey".
func (h *HomebrewPackageManager) Name() string {
	return "homebrew"
}

// Order implements PackageManager.Order by returning 2.
func (h *HomebrewPackageManager) Order() int {
	return 3
}

// IsSupported implements PackageManager.IsSupported by returning false.
func (h *HomebrewPackageManager) IsSupported() bool {
	// TODO: Implement this function.
	return false
}

// IsInstalled implements PackageManager.IsInstalled by returning an error.
func (h *HomebrewPackageManager) IsInstalled() (bool, error) {
	// TODO: Implement this function.
	return false, errors.New("not implemented")
}

// Install implements PackageManager.Install by returning an error.
func (h *HomebrewPackageManager) Install() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// Update implements PackageManager.Update by returning an error.
func (h *HomebrewPackageManager) Update() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// Uninstall implements PackageManager.Uninstall by returning an error.
func (h *HomebrewPackageManager) Uninstall() error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}

// InstalledPackages implements PackageManager.InstalledPackages by returning an error.
func (h *HomebrewPackageManager) InstalledPackages() ([]*Package, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// InstallPackage implements PackageManager.InstallPackage by returning an error.
func (h *HomebrewPackageManager) InstallPackage(packageName string, version *Version) (*Version, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// UpdatePackage implements PackageManager.UpdatePackage by returning an error.
func (h *HomebrewPackageManager) UpdatePackage(packageName string, version *Version) (*Version, error) {
	// TODO: Implement this function.
	return nil, errors.New("not implemented")
}

// UninstallPackage implements PackageManager.UninstallPackage by returning an error.
func (h *HomebrewPackageManager) UninstallPackage(packageName string) error {
	// TODO: Implement this function.
	return errors.New("not implemented")
}
