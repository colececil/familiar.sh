package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"strings"
)

// PackageCommand represents the "package" command.
type PackageCommand struct {
}

// Name returns the name of the command, as it appears on the command line while being used.
func (packageCommand *PackageCommand) Name() string {
	return "package"
}

// Description returns a short description of the command.
func (packageCommand *PackageCommand) Description() string {
	return "Manage packages for a given package manager."
}

// Documentation returns detailed documentation for the command.
func (packageCommand *PackageCommand) Documentation() string {
	return `The "package" command provides subcommands for adding, removing, and listing packages for a given package manager. It also allows you to specify the version of a package to install. It has the following subcommands:

  search
  info
  add
  remove
  update
  status
  import
`
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (packageCommand *PackageCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand must be included")
	}

	switch args[0] {
	case "add":
		subcommandArgs := args[1:]
		switch len(subcommandArgs) {
		case 1:
			return addPackageManager(subcommandArgs[0])
		case 2:
			return addPackage(subcommandArgs[0], subcommandArgs[1])
		default:
			return fmt.Errorf("wrong number of arguments")
		}
	case "remove":
		subcommandArgs := args[1:]
		switch len(subcommandArgs) {
		case 1:
			return removePackageManager(subcommandArgs[0])
		case 2:
			return removePackage(subcommandArgs[0], subcommandArgs[1])
		default:
			return fmt.Errorf("wrong number of arguments")
		}
	case "update":
		subcommandArgs := args[1:]
		switch len(subcommandArgs) {
		case 0:
			return updatePackages()
		case 1:
			return updatePackagesForPackageManager(subcommandArgs[0])
		case 2:
			return updatePackage(subcommandArgs[0], subcommandArgs[1])
		default:
			return fmt.Errorf("wrong number of arguments")
		}
	case "status":
		subcommandArgs := args[1:]
		switch len(subcommandArgs) {
		case 0:
			return getStatus()
		case 1:
			return getStatusForPackageManager(subcommandArgs[0])
		case 2:
			return getStatusForPackage(subcommandArgs[0], subcommandArgs[1])
		default:
			return fmt.Errorf("wrong number of arguments")
		}
	case "import":
		subcommandArgs := args[1:]
		switch len(subcommandArgs) {
		case 0:
			return importPackages()
		case 1:
			return importPackagesFromPackageManager(subcommandArgs[0])
		default:
			return fmt.Errorf("wrong number of arguments")
		}
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}

// addPackageManager adds the package manager of the given name to the config file.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to add.
func addPackageManager(packageManagerName string) error {
	if err := config.AddPackageManager(packageManagerName); err != nil {
		return err
	}

	fmt.Println("Package manager added.")
	return nil
}

// removePackageManager removes the package manager of the given name from the config file.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to remove.
func removePackageManager(packageManagerName string) error {
	if err := config.RemovePackageManager(packageManagerName); err != nil {
		return err
	}

	fmt.Println("Package manager removed.")
	return nil
}

// addPackage adds the given package to the config file under the given package manager. After that, it installs the
// package using the package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to use.
//   - packageName: The name of the package to add.
func addPackage(packageManagerName string, packageName string) error {
	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.Update(); err != nil {
		return err
	}

	installedVersion, err := packageManager.InstallPackage(packageName, nil)
	if err != nil {
		return err
	}

	return config.AddPackage(packageManagerName, packageName, installedVersion)
}

// removePackage removes the given package from the config file under the given package manager. After that, it
// uninstalls the package using the package manager.
//
// It takes the following parameters:
//   - packageManagerName: The name of the package manager to use.
//   - packageName: The name of the package to remove.
func removePackage(packageManagerName string, packageName string) error {
	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.UninstallPackage(packageName); err != nil {
		return err
	}

	return config.RemovePackage(packageManagerName, packageName)
}

// updatePackages updates all currently installed packages for all package managers that are both supported and
// installed.
func updatePackages() error {
	packageManagers := packagemanagers.GetAllPackageManagers()

	for _, packageManager := range packageManagers {
		if isSupported := packageManager.IsSupported(); !isSupported {
			continue
		}

		isInstalled, err := packageManager.IsInstalled()
		if err != nil {
			return err
		}

		if isInstalled {
			if err := updatePackagesForPackageManager(packageManager.Name()); err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipping package manager \"%s\" because it is not installed.\n", packageManager.Name())
		}
	}

	return nil
}

// updatePackagesForPackageManager updates all currently installed packages for the package manager of the given name.
func updatePackagesForPackageManager(packageManagerName string) error {
	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	configContents, err := config.ReadConfigFile()
	if err != nil {
		return err
	}

	configuredPackages := make(map[string]*packagemanagers.Version)
	for _, configuredPackageManager := range configContents.PackageManagers {
		if configuredPackageManager.Name == packageManagerName {
			for _, configuredPackage := range configuredPackageManager.Packages {
				configuredPackages[configuredPackage.Name] =
					&packagemanagers.Version{VersionString: configuredPackage.Version}
			}
			break
		}
	}

	if err := packageManager.Update(); err != nil {
		return err
	}

	installedPackages, err := packageManager.InstalledPackages()
	if err != nil {
		return err
	}

	for _, installedPackage := range installedPackages {
		if installedPackage.LatestVersion.IsGreaterThan(installedPackage.InstalledVersion) {
			newVersion, err := packageManager.UpdatePackage(installedPackage.Name, nil)
			if err != nil {
				return err
			}

			if configuredPackages[installedPackage.Name] != nil &&
				configuredPackages[installedPackage.Name].IsLessThan(newVersion) {
				if err := config.UpdatePackage(packageManagerName, installedPackage.Name, newVersion); err != nil {
					return err
				}
			}
		} else {
			fmt.Printf("Skipping package \"%s\" because it is already up to date.\n", installedPackage.Name)
		}
	}

	return nil
}

// updatePackage updates the given package for the given package manager.
func updatePackage(packageManagerName string, packageName string) error {
	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.Update(); err != nil {
		return err
	}

	installedPackages, err := packageManager.InstalledPackages()
	if err != nil {
		return err
	}

	for _, installedPackage := range installedPackages {
		if installedPackage.Name == packageName {
			if installedPackage.LatestVersion.IsGreaterThan(installedPackage.InstalledVersion) {
				newVersion, err := packageManager.UpdatePackage(packageName, nil)
				if err != nil {
					return err
				}

				configContents, err := config.ReadConfigFile()
				if err != nil {
					return err
				}

				for _, configuredPackageManager := range configContents.PackageManagers {
					if configuredPackageManager.Name == packageManagerName {
						for _, configuredPackage := range configuredPackageManager.Packages {
							if configuredPackage.Name == packageName {
								configuredVersion := &packagemanagers.Version{VersionString: configuredPackage.Version}
								if configuredVersion.IsLessThan(newVersion) {
									err := config.UpdatePackage(packageManagerName, packageName, newVersion)
									if err != nil {
										return err
									}
								}
								break
							}
						}
						break
					}
				}
			} else {
				fmt.Printf("Package \"%s\" is already up to date.\n", packageName)
			}

			return nil
		}
	}

	fmt.Printf("Package \"%s\" is not installed.\n", packageName)
	return nil
}

// getStatus prints the status for all package managers supported on the current machine.
func getStatus() error {
	packageManagers := packagemanagers.GetAllPackageManagers()

	for _, packageManager := range packageManagers {
		if isSupported := packageManager.IsSupported(); !isSupported {
			continue
		}

		isInstalled, err := packageManager.IsInstalled()
		if err != nil {
			return err
		}

		if isInstalled {
			if err := getStatusForPackageManager(packageManager.Name()); err != nil {
				return err
			}
		} else {
			fmt.Printf("Package manager \"%s\" is not installed.\n", packageManager.Name())
		}
	}

	return nil
}

// getStatusForPackageManager prints the status for the package manager of the given name.
func getStatusForPackageManager(packageManagerName string) error {
	configContents, err := config.ReadConfigFile()
	if err != nil {
		return err
	}

	configuredPackages := make(map[string]config.ConfiguredPackage)
	for _, configuredPackageManager := range configContents.PackageManagers {
		if configuredPackageManager.Name == packageManagerName {
			for _, configuredPackage := range configuredPackageManager.Packages {
				configuredPackages[configuredPackage.Name] = configuredPackage
			}
			break
		}
	}

	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.Update(); err != nil {
		return err
	}

	installedPackagesSlice, err := packageManager.InstalledPackages()
	if err != nil {
		return err
	}

	installedPackages := make(map[string]*packagemanagers.Package)
	for _, installedPackage := range installedPackagesSlice {
		installedPackages[installedPackage.Name] = installedPackage
	}

	if len(configuredPackages) == 0 && len(installedPackages) == 0 {
		fmt.Printf("No packages configured or installed for package manager \"%s\".\n", packageManager.Name())
		return nil
	}

	fmt.Printf("Status of packages for package manager \"%s\":\n", packageManager.Name())
	var packageListStringBuilder strings.Builder

	for _, configuredPackage := range configuredPackages {
		packageListStringBuilder.WriteString(fmt.Sprintf("- %s\n", configuredPackage.Name))
		packageListStringBuilder.WriteString(fmt.Sprintf("  - Configured version: %s\n", configuredPackage.Version))

		packageListStringBuilder.WriteString("  - Installed version: ")
		installedPackage, installed := installedPackages[configuredPackage.Name]
		if installed {
			packageListStringBuilder.WriteString(fmt.Sprintf("%s\n", installedPackage.InstalledVersion))
		} else {
			packageListStringBuilder.WriteString("\n")
		}

		packageListStringBuilder.WriteString("  - Newer version: ")
		if installed && installedPackage.LatestVersion.IsGreaterThan(installedPackage.InstalledVersion) {
			packageListStringBuilder.WriteString(fmt.Sprintf("%s\n", installedPackage.LatestVersion))
		} else {
			packageListStringBuilder.WriteString("\n")
		}
	}

	for _, installedPackage := range installedPackages {
		if _, configured := configuredPackages[installedPackage.Name]; !configured {
			packageListStringBuilder.WriteString(fmt.Sprintf("- %s\n", installedPackage.Name))
			packageListStringBuilder.WriteString("  - Configured version: \n")
			packageListStringBuilder.WriteString(fmt.Sprintf("  - Installed version: %s\n",
				installedPackage.InstalledVersion))

			packageListStringBuilder.WriteString("  - Newer version: ")
			if installedPackage.LatestVersion.IsGreaterThan(installedPackage.InstalledVersion) {
				packageListStringBuilder.WriteString(fmt.Sprintf("%s\n", installedPackage.LatestVersion))
			} else {
				packageListStringBuilder.WriteString("\n")
			}
		}
	}

	fmt.Print(packageListStringBuilder.String())
	return nil
}

// getStatusForPackage prints the status for the given package under the given package manager.
func getStatusForPackage(packageManagerName string, packageName string) error {
	configContents, err := config.ReadConfigFile()
	if err != nil {
		return err
	}

	var selectedPackageConfiguration config.ConfiguredPackage
	for _, configuredPackageManager := range configContents.PackageManagers {
		if configuredPackageManager.Name == packageManagerName {
			for _, configuredPackage := range configuredPackageManager.Packages {
				if configuredPackage.Name == packageName {
					selectedPackageConfiguration = configuredPackage
					break
				}
			}
		}
	}

	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.Update(); err != nil {
		return err
	}

	installedPackages, err := packageManager.InstalledPackages()
	if err != nil {
		return err
	}

	var selectedPackageInstallation *packagemanagers.Package
	for _, installedPackage := range installedPackages {
		if installedPackage.Name == packageName {
			selectedPackageInstallation = installedPackage
			break
		}
	}

	packageStringBuilder := strings.Builder{}
	packageStringBuilder.WriteString(fmt.Sprintf("Status of package \"%s\" for package manager \"%s\":\n", packageName,
		packageManager.Name()))

	packageStringBuilder.WriteString("- Configured version: ")
	if &selectedPackageConfiguration != nil {
		packageStringBuilder.WriteString(fmt.Sprintf("%s\n", selectedPackageConfiguration.Version))
	} else {
		packageStringBuilder.WriteString("\n")
	}

	packageStringBuilder.WriteString("- Installed version: ")
	if selectedPackageInstallation != nil {
		packageStringBuilder.WriteString(fmt.Sprintf("%s\n", selectedPackageInstallation.InstalledVersion))
	} else {
		packageStringBuilder.WriteString("\n")
	}

	packageStringBuilder.WriteString("- Newer version: ")
	if selectedPackageInstallation != nil &&
		selectedPackageInstallation.LatestVersion.IsGreaterThan(selectedPackageInstallation.InstalledVersion) {
		packageStringBuilder.WriteString(fmt.Sprintf("%s\n", selectedPackageInstallation.LatestVersion))
	} else {
		packageStringBuilder.WriteString("\n")
	}

	fmt.Print(packageStringBuilder.String())
	return nil
}

// importPackages imports all currently installed packages from all package managers that are both supported and
// installed.
func importPackages() error {
	packageManagers := packagemanagers.GetAllPackageManagers()

	for _, packageManager := range packageManagers {
		if isSupported := packageManager.IsSupported(); !isSupported {
			continue
		}

		isInstalled, err := packageManager.IsInstalled()
		if err != nil {
			return err
		}

		if isInstalled {
			if err := importPackagesFromPackageManager(packageManager.Name()); err != nil {
				return err
			}
		} else {
			fmt.Printf("Skipping package manager \"%s\" because it is not installed.\n", packageManager.Name())
		}
	}

	return nil
}

// importPackagesFromPackageManager imports all currently installed packages from the given package manager.
func importPackagesFromPackageManager(packageManagerName string) error {
	configContents, err := config.ReadConfigFile()
	if err != nil {
		return err
	}

	configuredPackageVersions := make(map[string]*packagemanagers.Version)
	for _, configuredPackageManager := range configContents.PackageManagers {
		if configuredPackageManager.Name == packageManagerName {
			for _, configuredPackage := range configuredPackageManager.Packages {
				configuredPackageVersions[configuredPackage.Name] =
					&packagemanagers.Version{VersionString: configuredPackage.Version}
			}
			break
		}
	}

	packageManager, err := packagemanagers.GetPackageManager(packageManagerName)
	if err != nil {
		return err
	}

	if err := packageManager.Update(); err != nil {
		return err
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

	packageManagerConfigUpdated := false
	for _, installedPackage := range installedPackages {
		configuredPackageVersion, isPresent := configuredPackageVersions[installedPackage.Name]
		if !isPresent {
			fmt.Printf("Adding package \"%s\" to configuration for package manager \"%s\".\n", installedPackage.Name,
				packageManagerName)

			err := config.AddPackage(packageManagerName, installedPackage.Name, installedPackage.InstalledVersion)
			if err != nil {
				return err
			}

			packageManagerConfigUpdated = true
		} else if configuredPackageVersion.IsGreaterThan(installedPackage.InstalledVersion) {
			fmt.Printf("Updating version of package \"%s\" in configuration for package manager \"%s\".\n",
				installedPackage.Name, packageManagerName)

			err := config.UpdatePackage(packageManagerName, installedPackage.Name, installedPackage.InstalledVersion)
			if err != nil {
				return err
			}

			packageManagerConfigUpdated = true
		}
	}

	if !packageManagerConfigUpdated {
		fmt.Printf("No packages to add or update in configuration for package manager \"%s\".\n", packageManagerName)
	}

	return nil
}
