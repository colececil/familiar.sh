package test

import "github.com/colececil/familiar.sh/internal/packagemanagers"

// PackageManagerDouble is a test double that implements the PackageManager interface. All methods except for Name are
// stubbed out.
type PackageManagerDouble struct {
	packagemanagers.PackageManager
	name  string
	order int
}

// NewPackageManagerDouble returns a new instance of PackageManagerDouble.
func NewPackageManagerDouble(name string, order int) *PackageManagerDouble {
	return &PackageManagerDouble{
		name:  name,
		order: order,
	}
}

// Name returns the name of the package manager.
func (d *PackageManagerDouble) Name() string {
	return d.name
}

// Order returns the order of the package manager.
func (d *PackageManagerDouble) Order() int {
	return d.order
}
