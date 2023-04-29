package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
)

type AttuneCommand struct {
}

// Name returns the name of the command, as it appears on the command line while being used.
func (attuneCommand *AttuneCommand) Name() string {
	return "attune"
}

// Description returns a short description of the command.
func (attuneCommand *AttuneCommand) Description() string {
	return "Sync the current machine with the config file."
}

// Documentation returns detailed documentation for the command.
func (attuneCommand *AttuneCommand) Documentation() string {
	return "Set up the current machine so it matches the shared configuration. To do this, Familiar.sh will perform " +
		"the following operations as needed: installing packages, uninstalling packages, copying files, and running " +
		"scripts."
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (attuneCommand *AttuneCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("the \"attune\" command does not take any arguments")
	}

	configContents, err := config.ReadConfigFile()
	if err != nil {
		return err
	}

	for _, packageManagerInConfig := range configContents.PackageManagers {
		packageManager, err := packagemanagers.GetPackageManager(packageManagerInConfig.Name)
		if err != nil {
			return err
		}

		installed, err := packageManager.IsInstalled()
		if err != nil {
			return err
		}

		if !installed {
			if err := packageManager.Install(); err != nil {
				return err
			}
		} else {
			fmt.Printf("Package manager \"%s\" is already installed.\n", packageManager.Name())
		}

		installedPackages, err := packageManager.InstalledPackages()
		if err != nil {
			return err
		}

		if len(installedPackages) > 0 {
			fmt.Printf("Packages currently installed with %s:\n", packageManager.Name())
			for _, installedPackage := range installedPackages {
				fmt.Printf("- %s, version %s\n", installedPackage.Name, installedPackage.Version)
			}
		} else {
			fmt.Printf("No packages are currently installed with %s.\n", packageManager.Name())
		}

		if len(packageManagerInConfig.Packages) > 0 {
			fmt.Printf("Packages configured to be installed with %s:\n", packageManager.Name())
			for _, packageInConfig := range packageManagerInConfig.Packages {
				fmt.Printf("- %s, version %s\n", packageInConfig.Name, packageInConfig.Version)
			}
		}

		var desiredPackageVersions = make(map[string]*packagemanagers.Version)
		for _, packageInConfig := range packageManagerInConfig.Packages {
			desiredPackageVersions[packageInConfig.Name] =
				&packagemanagers.Version{VersionString: packageInConfig.Version}
		}

		var installedPackageVersions = make(map[string]*packagemanagers.Version)
		for _, installedPackage := range installedPackages {
			installedPackageVersions[installedPackage.Name] = installedPackage.Version
		}

		// Update packages that are installed but have a lower version than the one in the config file.
		for packageName, desiredPackageVersion := range desiredPackageVersions {
			if installedPackageVersion, isPresent := installedPackageVersions[packageName]; isPresent {
				if installedPackageVersion.IsLessThan(desiredPackageVersion) {
					newVersion, err := packageManager.UpdatePackage(packageName, desiredPackageVersion)
					if err != nil {
						return err
					}

					if newVersion.IsGreaterThan(desiredPackageVersion) {
						err = config.UpdatePackage(packageManager.Name(), packageName, newVersion)
						if err != nil {
							return err
						}
					}
				}
			}
		}

		// Install packages that are in the config file but not installed.
		for packageName, desiredPackageVersion := range desiredPackageVersions {
			if _, isPresent := installedPackageVersions[packageName]; !isPresent {
				newVersion, err := packageManager.InstallPackage(packageName, desiredPackageVersion)
				if err != nil {
					return err
				}

				if newVersion.IsGreaterThan(desiredPackageVersion) {
					err = config.UpdatePackage(packageManager.Name(), packageName, newVersion)
					if err != nil {
						return err
					}
				}
			}
		}

		// Uninstall packages that are installed but not in the config file.
		for packageName := range installedPackageVersions {
			if _, isPresent := desiredPackageVersions[packageName]; !isPresent {
				if err := packageManager.UninstallPackage(packageName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
