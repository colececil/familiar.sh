//go:build wireinject

package main

import (
	"github.com/adrg/xdg"
	"github.com/colececil/familiar.sh/internal/commands"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/system"
	"github.com/google/wire"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// This value is overridden at build time using `-ldflags`.
const projectVersion = "0.0.0"

var providers = wire.NewSet(
	getFamiliarVersion,
	getCurrentOperatingSystem,
	getCommands,
	getPackageManagers,
	getCreateShellCommandFunc,
	getCreateShellCommandRunnerFunc,
	getXdgConfigHomeGetter,
	getAbsPathConverter,
	getDirPathGetter,
	getDirCreator,
	getFileExistenceChecker,
	getFileReader,
	getFileCreator,
	getFileSystem,
	getOutputWriter,
	commands.NewCommandRegistry,
	commands.NewVersionCommand,
	commands.NewAttuneCommand,
	commands.NewConfigCommand,
	commands.NewPackageCommand,
	commands.NewHelpCommand,
	config.NewConfigService,
	packagemanagers.NewPackageManagerRegistry,
	packagemanagers.NewScoopPackageManager,
	packagemanagers.NewChocolateyPackageManager,
	packagemanagers.NewHomebrewPackageManager,
	system.NewOperatingSystemService,
)

// InitializeCommandRegistry tells Wire how to create an injector for CommandRegistry.
func InitializeCommandRegistry() commands.CommandRegistry {
	wire.Build(providers)
	return commands.CommandRegistry{}
}

// getFamiliarVersion returns the version of Familiar.sh.
func getFamiliarVersion() commands.FamiliarVersionString {
	version := commands.FamiliarVersionString(projectVersion)
	return version
}

// getCommands returns a slice containing all commands.
func getCommands(
	helpCommand *commands.HelpCommand,
	versionCommand *commands.VersionCommand,
	attuneCommand *commands.AttuneCommand,
	configCommand *commands.ConfigCommand,
	packageCommand *commands.PackageCommand,
) []commands.Command {
	return []commands.Command{
		helpCommand,
		versionCommand,
		attuneCommand,
		configCommand,
		packageCommand,
	}
}

// getPackageManagers returns a slice containing all package managers.
func getPackageManagers(
	scoopPackageManager *packagemanagers.ScoopPackageManager,
	chocolateyPackageManager *packagemanagers.ChocolateyPackageManager,
	homebrewPackageManager *packagemanagers.HomebrewPackageManager,
) []packagemanagers.PackageManager {
	return []packagemanagers.PackageManager{
		scoopPackageManager,
		chocolateyPackageManager,
		homebrewPackageManager,
	}
}

// getCreateShellCommandFunc returns system.NewShellCommand as a system.CreateShellCommandFunc.
func getCreateShellCommandFunc() system.CreateShellCommandFunc {
	return system.NewShellCommand
}

// getCreateShellCommandRunnerFunc returns system.NewShellCommandRunner as a system.CreateShellCommandRunnerFunc.
func getCreateShellCommandRunnerFunc() system.CreateShellCommandRunnerFunc {
	return system.NewShellCommandRunner
}

// getCurrentOperatingSystem returns the current operating system.
func getCurrentOperatingSystem() system.OperatingSystem {
	return system.OperatingSystem(runtime.GOOS)
}

// getXdgConfigHomeGetter returns a function the returns the XDG config home directory as a system.XdgConfigHomeGetter.
func getXdgConfigHomeGetter() system.XdgConfigHomeGetter {
	return system.XdgConfigHomeGetterFunc(
		func() string {
			return xdg.ConfigHome
		})
}

// getAbsPathConverter returns filepath.Abs as a system.AbsPathConverter.
func getAbsPathConverter() system.AbsPathConverter {
	return system.AbsPathConverterFunc(filepath.Abs)
}

// getDirPathGetter returns filepath.Dir as a system.PathDirGetter.
func getDirPathGetter() system.PathDirGetter {
	return system.PathDirGetterFunc(filepath.Dir)
}

// getFileExtensionGetter returns filepath.Ext as a system.FileExtensionGetter.
func getFileExtensionGetter() system.FileExtensionGetter {
	return system.FileExtensionGetterFunc(filepath.Ext)
}

// getDirCreator returns os.MkdirAll as a system.DirCreator.
func getDirCreator() system.DirCreator {
	return system.DirCreatorFunc(os.MkdirAll)
}

// getFileExistenceChecker returns system.FileExists as a system.FileExistenceChecker.
func getFileExistenceChecker() system.FileExistenceChecker {
	return system.FileExistenceCheckerFunc(system.FileExists)
}

// getFileReader returns os.ReadFile as a system.FileReader.
func getFileReader() system.FileReader {
	return system.FileReaderFunc(os.ReadFile)
}

// getFileCreator returns os.Create as a system.FileCreator.
func getFileCreator() system.FileCreator {
	return system.FileCreatorFunc(
		func(path string) (io.WriteCloser, error) {
			file, err := os.Create(path)
			return io.WriteCloser(file), err
		})
}

// getFileSystem returns an instance of config.FileSystem.
func getFileSystem() config.FileSystem {
	return config.FileSystem(
		struct {
			system.XdgConfigHomeGetter
			system.AbsPathConverter
			system.PathDirGetter
			system.FileExtensionGetter
			system.DirCreator
			system.FileExistenceChecker
			system.FileReader
			system.FileCreator
		}{
			getXdgConfigHomeGetter(),
			getAbsPathConverter(),
			getDirPathGetter(),
			getFileExtensionGetter(),
			getDirCreator(),
			getFileExistenceChecker(),
			getFileReader(),
			getFileCreator(),
		})
}

// getOutputWriter returns os.Stdout as an io.Writer.
func getOutputWriter() io.Writer {
	return os.Stdout
}
