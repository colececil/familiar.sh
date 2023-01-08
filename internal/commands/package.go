package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
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

  status
  search
  info
  add
  remove
  update`
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
