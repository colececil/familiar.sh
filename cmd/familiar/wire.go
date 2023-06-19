//go:build wireinject

package main

import (
	"github.com/colececil/familiar.sh/internal/commands"
	"github.com/colececil/familiar.sh/internal/config"
	"github.com/colececil/familiar.sh/internal/packagemanagers"
	"github.com/colececil/familiar.sh/internal/system"
	"github.com/google/wire"
)

// This value is overridden at build time using `-ldflags`.
const projectVersion = "0.0.0"

var providers = wire.NewSet(
	getFamiliarVersion,
	commands.NewCommandRegistry,
	commands.NewVersionCommand,
	commands.NewAttuneCommand,
	commands.NewConfigCommand,
	commands.NewPackageCommand,
	commands.NewHelpCommand,
	config.NewConfigService,
	packagemanagers.NewPackageManagerRegistry,
	packagemanagers.NewScoopPackageManager,
	system.NewIsWindowsFunc,
	system.NewOperatingSystemService,
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
