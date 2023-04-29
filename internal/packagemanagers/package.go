package packagemanagers

// Package represents a package that is managed by a package manager.
type Package struct {
	Name    string
	Version *Version
}
