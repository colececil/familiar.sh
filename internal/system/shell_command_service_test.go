package system_test

import (
	"bytes"
	. "github.com/colececil/familiar.sh/internal/system"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"time"
)

var _ = Describe("ShellCommandService", func() {
	var shellCommandDouble *test.ShellCommandDouble
	var outputWriterDouble *bytes.Buffer
	var shellCommandService *ShellCommandService

	const expectedProgram = "program"
	const expectedProgramArg = "arg"
	const programStdout = `This is
the program's standard output
`
	const programStderr = `This is
the program's error output
`
	const programExitCode = 0

	BeforeEach(func() {
		createShellCommandFuncDouble := func(program string, args ...string) ShellCommand {
			if program == expectedProgram && len(args) == 1 && args[0] == expectedProgramArg {
				return shellCommandDouble
			}
			return NewRealShellCommand("")
		}
		outputWriterDouble = new(bytes.Buffer)
		shellCommandService = NewShellCommandService(
			createShellCommandFuncDouble,
			NewRunShellCommandFunc(),
			outputWriterDouble,
		)
	})

	When("`printOutput` is set to `true`", func() {
		It("should print the command output", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					completionChannel chan<- int) {

					if _, err := stdoutWriter.Write([]byte(programStdout)); err != nil {
						errorChannel <- err
						return
					}
					completionChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, true, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(programStdout))
		})

		It("should print both the stdout and stderr command output", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					completionChannel chan<- int) {

					if _, err := stdoutWriter.Write([]byte(programStdout)); err != nil {
						errorChannel <- err
						return
					}
					time.Sleep(100 * time.Millisecond)
					if _, err := stderrWriter.Write([]byte(programStderr)); err != nil {
						errorChannel <- err
						return
					}
					completionChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, true, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(programStdout + programStderr))
		})

		It("should print each line of the command output as soon as it becomes available", func() {
		})
	})

	When("`printOutput` is set to `false`", func() {
		It("should not print the command output", func() {
		})
	})

	When("`resultCaptureRegex` contains a regex", func() {
		It("should return the first match of the regex from the cumulative command output, and the returned result "+
			"should come from the regex's first capturing group", func() {
		})

		It("should only return a result after the command finishes running", func() {
		})

		It("should return an empty string if there are no regex matches in the cumulative command output", func() {
		})

		It("should only consider the stdout command output when finding a regex match", func() {
		})

		It("should return an empty string if the regex contains no capturing group", func() {
		})
	})

	When("`resultCaptureRegex` is set to `nil`", func() {
		It("should return an empty string", func() {
		})
	})

	It("should return an error if the command fails to start", func() {
	})

	It("should return an error if the command fails to return its stdout pipe", func() {
	})

	It("should return an error if the command fails to return its stderr pipe", func() {
	})

	It("should return an error if the command returns a non-zero exit code", func() {
	})

	It("should return an error if the command fails to run", func() {
	})

	It("should return an error if there is an issue reading the stdout command output", func() {
	})

	It("should return an error if there is an issue reading the stderr command output", func() {
	})
})
