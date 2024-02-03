package commands_test

import (
	. "github.com/colececil/familiar.sh/internal/commands"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommandRegistry", func() {
	var commandRegistry CommandRegistry

	BeforeEach(func() {
		commandRegistry = NewCommandRegistry(
			test.NewCommandDouble("package", 5),
			test.NewCommandDouble("config", 4),
			test.NewCommandDouble("attune", 3),
			test.NewCommandDouble("version", 2),
			test.NewCommandDouble("help", 1),
		)
	})

	Describe("NewCommandRegistry", func() {
		const panicMessage = "command registry does not contain the expected commands"

		It("should panic if the command registry does not contain the expected commands", func() {
			Expect(func() {
				NewCommandRegistry(
					test.NewCommandDouble("help", 1),
					test.NewCommandDouble("version", 2),
					test.NewCommandDouble("attune", 3),
					test.NewCommandDouble("config", 4),
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are less than 1", func() {
			Expect(func() {
				NewCommandRegistry(
					test.NewCommandDouble("help", 0),
					test.NewCommandDouble("version", 2),
					test.NewCommandDouble("attune", 3),
					test.NewCommandDouble("config", 4),
					test.NewCommandDouble("package", 5),
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are greater than the number of commands", func() {
			Expect(func() {
				NewCommandRegistry(
					test.NewCommandDouble("help", 1),
					test.NewCommandDouble("version", 2),
					test.NewCommandDouble("attune", 3),
					test.NewCommandDouble("config", 4),
					test.NewCommandDouble("package", 6),
				)
			}).To(PanicWith(panicMessage))
		})

		It("should panic if any of the commands' orders are not unique", func() {
			Expect(func() {
				NewCommandRegistry(
					test.NewCommandDouble("help", 1),
					test.NewCommandDouble("version", 2),
					test.NewCommandDouble("attune", 3),
					test.NewCommandDouble("config", 3),
					test.NewCommandDouble("package", 5),
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
