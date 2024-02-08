package commands_test

import (
	. "github.com/colececil/familiar.sh/internal/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ovechkin-dm/mockio/mock"
)

var _ = Describe("CommandRegistry", func() {
	var commandRegistry CommandRegistry
	var helpCommand Command
	var versionCommand Command
	var attuneCommand Command
	var configCommand Command
	var packageCommand Command

	BeforeEach(func() {
		mock.SetUp(GinkgoT())

		helpCommand = mock.Mock[Command]()
		mock.WhenSingle(helpCommand.Name()).ThenReturn("help")
		mock.WhenSingle(helpCommand.Order()).ThenReturn(1)

		versionCommand = mock.Mock[Command]()
		mock.WhenSingle(versionCommand.Name()).ThenReturn("version")
		mock.WhenSingle(versionCommand.Order()).ThenReturn(2)

		attuneCommand = mock.Mock[Command]()
		mock.WhenSingle(attuneCommand.Name()).ThenReturn("attune")
		mock.WhenSingle(attuneCommand.Order()).ThenReturn(3)

		configCommand = mock.Mock[Command]()
		mock.WhenSingle(configCommand.Name()).ThenReturn("config")
		mock.WhenSingle(configCommand.Order()).ThenReturn(4)

		packageCommand = mock.Mock[Command]()
		mock.WhenSingle(packageCommand.Name()).ThenReturn("package")
		mock.WhenSingle(packageCommand.Order()).ThenReturn(5)

		commandRegistry = NewCommandRegistry(
			packageCommand,
			configCommand,
			attuneCommand,
			versionCommand,
			helpCommand,
		)
	})

	Describe("NewCommandRegistry", func() {
		const panicMessage = "command registry does not contain the expected commands"

		It("should panic if the command registry does not contain the expected commands", func() {
			Expect(func() {
				NewCommandRegistry(
					helpCommand,
					versionCommand,
					attuneCommand,
					configCommand,
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are less than 1", func() {
			mock.SetUp(GinkgoT())

			helpCommand = mock.Mock[Command]()
			mock.WhenSingle(helpCommand.Name()).ThenReturn("help")
			mock.WhenSingle(helpCommand.Order()).ThenReturn(0)

			Expect(func() {
				NewCommandRegistry(
					helpCommand,
					versionCommand,
					attuneCommand,
					configCommand,
					packageCommand,
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are greater than the number of commands", func() {
			mock.SetUp(GinkgoT())

			packageCommand = mock.Mock[Command]()
			mock.WhenSingle(packageCommand.Name()).ThenReturn("package")
			mock.WhenSingle(packageCommand.Order()).ThenReturn(6)

			Expect(func() {
				NewCommandRegistry(
					helpCommand,
					versionCommand,
					attuneCommand,
					configCommand,
					packageCommand,
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are not unique", func() {
			mock.SetUp(GinkgoT())

			configCommand = mock.Mock[Command]()
			mock.WhenSingle(configCommand.Name()).ThenReturn("config")
			mock.WhenSingle(configCommand.Order()).ThenReturn(3)

			Expect(func() {
				NewCommandRegistry(
					helpCommand,
					versionCommand,
					attuneCommand,
					configCommand,
					packageCommand,
				)
			}).To(PanicWith(panicMessage))
		})
	})

	Describe("GetAllCommands", func() {
		It("should return a slice containing all commands, sorted by their returned orders", func() {
			result := commandRegistry.GetAllCommands()

			Expect(result).To(HaveLen(5))
			Expect(result[0].Name()).To(Equal("help"))
			Expect(result[1].Name()).To(Equal("version"))
			Expect(result[2].Name()).To(Equal("attune"))
			Expect(result[3].Name()).To(Equal("config"))
			Expect(result[4].Name()).To(Equal("package"))
		})
	})

	Describe("GetCommand", func() {
		It("should return the command of the given name if it is in the registry", func() {
			command, err := commandRegistry.GetCommand("package")
			Expect(err).To(BeNil())
			Expect(command.Name()).To(Equal("package"))
		})

		It("should return an error if the command with the given name is not in the registry", func() {
			_, err := commandRegistry.GetCommand("invalid")
			Expect(err.Error()).To(Equal("command not valid"))
		})
	})
})
