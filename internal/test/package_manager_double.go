package test

import "github.com/colececil/familiar.sh/internal/packagemanagers"

// PackageManagerDouble is a test double that implements the PackageManager interface. All methods except for Name are
// stubbed out.
type PackageManagerDouble struct {
	packagemanagers.PackageManager
	name string
}

// NewPackageManagerDouble returns a new instance of PackageManagerDouble.
func NewPackageManagerDouble(name string) *PackageManagerDouble {
	return &PackageManagerDouble{
		name: name,
	}
}

// Name returns the name of the package manager.
func (packageManagerDouble *PackageManagerDouble) Name() string {
	return packageManagerDouble.name
}
