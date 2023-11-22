package test

import "github.com/colececil/familiar.sh/internal/packagemanagers"

// PackageManagerDouble is a test double that implements the PackageManager interface. All methods except for Name are
// stubbed out.
type PackageManagerDouble struct {
	name string
}

// NewPackageManagerDouble returns a new instance of PackageManagerDouble.
func NewPackageManagerDouble(name string) *PackageManagerDouble {
	return &PackageManagerDouble{
		name: name,
	}
}

func (packageManagerDouble *PackageManagerDouble) Name() string {
	return packageManagerDouble.name
}

func (packageManagerDouble *PackageManagerDouble) IsSupported() bool {
	return true
}

func (packageManagerDouble *PackageManagerDouble) IsInstalled() (bool, error) {
	return true, nil
}

func (packageManagerDouble *PackageManagerDouble) Install() error {
	return nil
}

func (packageManagerDouble *PackageManagerDouble) Update() error {
	return nil
}

func (packageManagerDouble *PackageManagerDouble) Uninstall() error {
	return nil
}

func (packageManagerDouble *PackageManagerDouble) InstalledPackages() ([]*packagemanagers.Package, error) {
	return []*packagemanagers.Package{}, nil
}

func (packageManagerDouble *PackageManagerDouble) InstallPackage(packageName string,
	version *packagemanagers.Version) (*packagemanagers.Version, error) {

	return nil, nil
}

func (packageManagerDouble *PackageManagerDouble) UpdatePackage(packageName string,
	version *packagemanagers.Version) (*packagemanagers.Version, error) {

	return nil, nil
}

func (packageManagerDouble *PackageManagerDouble) UninstallPackage(packageName string) error {
	return nil
}
