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

// NewPackageFromStrings creates a new instance of Package. This is an alternate constructor that takes in string
// representations of the versions.
func NewPackageFromStrings(name string, installedVersion string, latestVersion string) *Package {
	return &Package{
		Name:             name,
		InstalledVersion: NewVersion(installedVersion),
		LatestVersion:    NewVersion(latestVersion),
	}
}
