package packagemanagers

// PackageManager represents a package manager that is supported by Familiar.sh.
type PackageManager interface {
	// Name returns the name of the package manager.
	Name() string

	// IsSupported returns whether the package manager is supported on the current machine.
	IsSupported() bool

	// IsInstalled returns whether the package manager is installed.
	IsInstalled() (bool, error)

	// Install installs the package manager.
	Install() error

	// Update updates the package manager.
	Update() error

	// Uninstall uninstalls the package manager.
	Uninstall() error

	// InstalledPackages returns a slice containing information about all packages that are installed.
	InstalledPackages() ([]*Package, error)

	// InstallPackage installs the package of the given name. If a version is given, that specific version of the
	// package is installed. Otherwise, the latest version is installed.
	//
	// It takes the following parameters:
	//   - packageName: The name of the package to install.
	// 	 - version: The version of the package to install. If nil, the latest version is installed.
	//
	// It returns the version of the package that was installed.
	InstallPackage(packageName string, version *Version) (*Version, error)

	// UpdatePackage updates the package of the given name. If a version is given, that specific version of the package
	// is installed. Otherwise, the latest version is installed.
	//
	// It takes the following parameters:
	//   - packageName: The name of the package to install.
	// 	 - version: The version of the package to install. If nil, the latest version is installed.
	//
	// It returns the version of the package that was installed.
	UpdatePackage(packageName string, version *Version) (*Version, error)

	// UninstallPackage uninstalls the package of the given name.
	UninstallPackage(packageName string) error
}
