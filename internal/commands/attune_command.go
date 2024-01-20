package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
)

type AttuneCommand struct {
	configService          *config.ConfigService
	packageManagerRegistry packagemanagers.PackageManagerRegistry
}

// NewAttuneCommand creates a new instance of AttuneCommand.
func NewAttuneCommand(configService *config.ConfigService,
	packageManagerRegistry packagemanagers.PackageManagerRegistry) *AttuneCommand {
	return &AttuneCommand{
		configService:          configService,
		packageManagerRegistry: packageManagerRegistry,
	}
}

// Name returns the name of the command, as it appears on the command line while being used.
func (c *AttuneCommand) Name() string {
	return "attune"
}

// Order returns the order in which the command should be listed in the help command.
func (c *AttuneCommand) Order() int {
	return 3
}

// Description returns a short description of the command.
func (c *AttuneCommand) Description() string {
	return "Sync the current machine with the config file."
}

// Documentation returns detailed documentation for the command.
func (c *AttuneCommand) Documentation() string {
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
func (c *AttuneCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("the \"attune\" command does not take any arguments")
	}

	configContents, err := c.configService.GetConfig()
	if err != nil {
		return err
	}

	for _, packageManagerInConfig := range configContents.PackageManagers {
		packageManager, err := c.packageManagerRegistry.GetPackageManager(packageManagerInConfig.Name)
		if err != nil {
			return err
		}

		if isSupported := packageManager.IsSupported(); !isSupported {
			continue
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
			if err := packageManager.Update(); err != nil {
				return err
			}
		}

		installedPackages, err := packageManager.InstalledPackages()
		if err != nil {
			return err
		}

		if len(installedPackages) > 0 {
			fmt.Printf("Packages currently installed with %s:\n", packageManager.Name())
			for _, installedPackage := range installedPackages {
				fmt.Printf("- %s, version %s\n", installedPackage.Name, installedPackage.InstalledVersion)
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
			desiredPackageVersions[packageInConfig.Name] = packagemanagers.NewVersion(packageInConfig.Version)
		}

		var installedPackageVersions = make(map[string]*packagemanagers.Version)
		for _, installedPackage := range installedPackages {
			installedPackageVersions[installedPackage.Name] = installedPackage.InstalledVersion
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
						err = configContents.UpdatePackage(packageManager.Name(), packageName, newVersion)
						if err != nil {
							return err
						}

						if err = c.configService.SetConfig(configContents); err != nil {
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
					err = configContents.UpdatePackage(packageManager.Name(), packageName, newVersion)
					if err != nil {
						return err
					}

					if err = c.configService.SetConfig(configContents); err != nil {
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
