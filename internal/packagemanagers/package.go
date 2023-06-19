package packagemanagers

// Package represents a package that is managed by a package manager.
type Package struct {
	Name             string
	InstalledVersion *Version
	LatestVersion    *Version
}

// NewPackage creates a new instance of Package.
func NewPackage(name string, installedVersion *Version, latestVersion *Version) *Package {
	return &Package{
		Name:             name,
		InstalledVersion: installedVersion,
		LatestVersion:    latestVersion,
	}
}
