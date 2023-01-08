package packagemanagers

type PackageManager interface {
	// Name returns the name of the package manager.
	Name() string

	// InstallPackage installs the package of the given name. If a version is given, that specific version of the
	// package is installed. Otherwise, the latest version is installed.
	//
	// It returns the version of the package that was installed.
	InstallPackage(packageName string, version *string) (string, error)

	// UninstallPackage uninstalls the package of the given name.
	UninstallPackage(packageName string) error
}
