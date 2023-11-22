//go:build wireinject

package main

import (
	"github.com/colececil/familiar.sh/internal/commands"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/system"
	"github.com/google/wire"
	"io"
	"os"
	"runtime"
)

// This value is overridden at build time using `-ldflags`.
const projectVersion = "0.0.0"

var providers = wire.NewSet(
	getFamiliarVersion,
	getCurrentOperatingSystem,
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
	system.NewOperatingSystemService,
	system.NewCreateShellCommandFunc,
	system.NewRunShellCommandFunc,
	system.NewShellCommandService,
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

// getCurrentOperatingSystem returns the current operating system.
func getCurrentOperatingSystem() system.OperatingSystem {
	return system.OperatingSystem(runtime.GOOS)
}

// getOutputWriter returns the output writer.
func getOutputWriter() io.Writer {
	return os.Stdout
}
