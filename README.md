![Familiar.sh logo](logo/familiar-logo.png)

# Familiar.sh

A cross-platform CLI for setting up all your machines just the way you like and keeping them in sync. Supports Windows, MacOS, and Linux.

Do you work across multiple machines, but waste a lot of time trying to get them all set up with your preferred tools and configurations? Do you switch between different operating systems on a regular basis? Do you struggle to remember the slight inconsistencies in usage between various package managers? **Then Familiar.sh can help!**

Familiar.sh offers the following features:

- It provides a unified interface that serves as an abstraction over multiple commonly used package managers (Apt, Yum, Chocolatey, Scoop, Homebrew, SDKMan, etc.). As you install, update, and uninstall packages, Familiar.sh automatically tracks the changes in its configuration file.
- With the use of a cloud drive like Google Drive or Microsoft OneDrive, your Familiar.sh configuration can be shared and synced across multiple machines. Then just run `familiar attune` to get everything in order!
- If you place copies of your configuration files (.bashrc, .vimrc, SSH config, AWS config, etc.) next to your Familiar.sh configuration in your cloud drive, you can tell Familiar.sh where they belong and have it keep them in sync for you as well.
- You can even write scripts that perform custom machine setup and have Familiar.sh run them as needed.

## Usage

_Note: Familiar.sh is still early in development, and it is not yet usable. The commands are also subject to change until the CLI reaches a stable version._

### How to Install

_(Yet to be written.)_

### Commands

- **CLI Information**
  - `familiar help` (alias `--help`, `-h`): List help information. `help` can also be used to get information about individual subcommands (for example, you can get information about the `config` subcommand by running `familiar help config`).
  - `familiar version` (alias `--version`, `-v`): Print the installed version of Familiar.sh.
- **Shared Configuration**
  - `familiar attune` (alias `sync`): Set up the current machine so it matches the shared configuration. To do this, Familiar.sh will perform the following operations as needed: installing packages, uninstalling packages, copying files, and running scripts.
  - `familiar config`: Print the contents of the shared configuration file.
  - `familiar config location`: Print the config file location.
  - `familiar config location <path>`: Set the config file location to the given path.
- **Configuration Management**
  - `familiar file add <sourcePath> <destinationPath>`: Add the file at the given source path to the shared configuration, telling Familiar.sh it should be synced to the given destination path.
  - `familiar file remove <filename>`: Remove the given file from the shared configuration.
  - `familiar script add <path>`: Add the script at the given path to the shared configuration. The script will be run whenever `familiar attune` is run, so it should be idempotent.
    - Optional flags:
      - `--operating-systems <operatingSystems>`: Specify which operating systems the script should run on (by default, it runs on all operating systems). The operating systems should be a comma separated list - valid values are `windows`, `macos`, and `linux`. For example, `--operating-systems "macos, linux"` would specify that the script should only be run on MacOS and Linux.
      - `--preconditions <preconditions>`: Specify that the script should only be run when the given preconditions are met. _Note: The way of specifying preconditions is still being designed. More details to come._
  - `familiar script remove <filename>`: Remove the given script from the shared configuration.
- **Package Management**
  - **Package Search and Information**
    - `familiar package search <term>`: Search for packages with the given term under all installed package managers.
    - `familiar package search <packageManager> <term>`: Search for packages using the given term under the given package manager.
    - `familiar package info <packageManager> <term>`: Print information about the given package under the given package manager.
  - **Installation and Uninstallation**
    - `familiar package add <packageManager>` (alias `package install`): Install the given package manager.
    - `familiar package add <packageManager> <package>` (alias `package install`): Install the given package using the given package manager. This also adds the package to the shared configuration.
      - Optional flags:
        - `--no-save`: Perform the operation without updating the shared configuration.
    - `familiar package remove <packageManager>` (alias `package uninstall`): Uninstall the given package manager, along with all its installed packages.
    - `familiar package remove <packageManager> <package>` (alias `package uninstall`): Uninstall the given package using the given package manager. This also removes the package from the shared configuration.
      - Optional flags:
        - `--no-save`: Perform the operation without updating the shared configuration.
  - **Updating**
    - `familiar package update` (alias `package upgrade`): Update all installed packages to the latest available version. This also updates the package versions in the shared configuration.
      - Optional flags:
        - `--no-save`: Perform the operation without updating the shared configuration.
    - `familiar package update <packageManager>` (alias `package upgrade`): Update all installed packages under the given package manager to the latest available version. This also updates the package versions in the shared configuration.
      - Optional flags:
        - `--no-save`: Perform the operation without updating the shared configuration.
    - `familiar package update <packageManager> <package>` (alias `package upgrade`): Update the given package under the given package manager to the latest available version. This also updates the package version in the shared configuration.
      - Optional flags:
        - `--no-save`: Perform the operation without updating the shared configuration.
  - **Status of Installation and Updates**
    - `familiar package status`: Show the status of all configured/installed packages for all installed package managers, along with any available updates.
    - `familiar package status <packageManager>`: Show the status of all configured/installed packages under the given package manager, along with any available updates.
    - `familiar package status <packageManager> <package>`: Show the status for the given package under the given package manager.
  - **Importing**
    - `familiar package import`: For all package managers that are supported and installed, import all currently installed packages into the shared configuration. This is helpful for getting started with Familiar.sh on a machine that already has a lot of packages installed.
    - `familiar package import <packageManager>`: For the given package manager, import all currently installed packages into the shared configuration. This is helpful for getting started with Familiar.sh on a machine that already has a lot of packages installed.

### Examples

_(Yet to be written.)_

## Development

### Prerequisites

In order to develop Familiar.sh on your local machine, you'll first need to follow these steps:

1. Make sure the following are installed:
    - [Go v1.20](https://go.dev/doc/install): Go was chosen as the programming language for this project because of its ability to easily create executables for different operating systems, and because the executables it creates require no additional dependencies.
    - [Wire CLI v0.5.0](https://github.com/google/wire): This project uses Wire, which is a compile-time dependency injection framework. The Wire CLI package (located at <https://github.com/google/wire/tree/main/cmd/wire>) is needed in order to generate the injection code before running or building the project.
2. Download the project dependencies by running `go mod download` from the project root.
3. Generate the dependency injection code by running `wire gen ./...` from the project root. (This generated code will be saved at `cmd/familiar/wire_gen.go`.)

### Running the Project

To run the project, run the following from the project root, replacing `<command>` with the command you want to run:

```
go run ./cmd/familiar/ <command>
```

### Building the Project

To build the project, run the following from the project root:

```
go build ./cmd/familiar/
```

### Passing in the Project Version

The project version can be passed in when building or running the project using `-ldflags '-X main.projectVersion=<version>'`, replacing `<version>` with the desired version. For example:

```
go build -ldflags '-X main.projectVersion=1.0.0' ./cmd/familiar/
```

This causes the `projectVersion` constant in `cmd/familiar/wire_gen.go` to be set to the given version in the executable, rather than using the default value of `"0.0.0"`.
