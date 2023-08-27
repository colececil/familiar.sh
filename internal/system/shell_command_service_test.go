package system_test

import (
	"bytes"
	. "github.com/colececil/familiar.sh/internal/system"
	"github.com/colececil/familiar.sh/internal/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"regexp"
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
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(programStdout))
					exitCodeChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, true, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(programStdout))
		})

		It("should print both the stdout and stderr command output", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(programStdout))
					time.Sleep(100 * time.Millisecond)
					_, _ = stderrWriter.Write([]byte(programStderr))
					exitCodeChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, true, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(programStdout + programStderr))
		})

		It("should print each line of the command output as soon as it becomes available", func() {
			stdoutLine1 := "stdout line 1"
			stdoutLine2 := "stdout line 2"
			stdoutLine3 := "stdout line 3"
			stderrLine1 := "stderr line 1"
			stderrLine2 := "stderr line 2"

			outputChannel := make(chan string, 5)
			testCompletionChannel := make(chan bool)

			go func() {
				defer GinkgoRecover()

				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n"))
				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n"))
				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n" + stderrLine1 + "\n"))
				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n" + stderrLine1 + "\n" + stdoutLine2 + "\n"))
				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n" + stderrLine1 + "\n" + stdoutLine2 + "\n" +
					stdoutLine3 + "\n"))
				Expect(<-outputChannel).To(Equal(stdoutLine1 + "\n" + stderrLine1 + "\n" + stdoutLine2 + "\n" +
					stdoutLine3 + "\n"))

				testCompletionChannel <- true
			}()

			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(stdoutLine1 + "\n"))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					_, _ = stderrWriter.Write([]byte(stderrLine1))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					_, _ = stderrWriter.Write([]byte("\n"))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					_, _ = stdoutWriter.Write([]byte(stdoutLine2 + "\n"))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					_, _ = stdoutWriter.Write([]byte(stdoutLine3 + "\n"))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					_, _ = stdoutWriter.Write([]byte(stderrLine2))
					time.Sleep(100 * time.Millisecond)
					outputChannel <- outputWriterDouble.String()

					exitCodeChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, true, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(stdoutLine1 + "\n" + stderrLine1 + "\n" + stdoutLine2 + "\n" +
				stdoutLine3 + "\n" + stderrLine2 + "\n"))
			Eventually(testCompletionChannel).Should(Receive())
		})
	})

	When("`printOutput` is set to `false`", func() {
		It("should not print the command output", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(programStdout))
					_, _ = stderrWriter.Write([]byte(programStderr))
					exitCodeChannel <- programExitCode
				},
			)

			_, err := shellCommandService.RunShellCommand(expectedProgram, false, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(outputWriterDouble.String()).To(Equal(""))
		})
	})

	When("`resultCaptureRegex` contains a regex", func() {
		var output string
		var regex *regexp.Regexp

		BeforeEach(func() {
			output = "Line 1\nLine 2\nLine 3\nThe result is 123\n"
			regex, _ = regexp.Compile("The result is (\\d+)")
		})

		It("should return the first match of the regex from the cumulative command output, and the returned result "+
			"should come from the regex's first capturing group", func() {

			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(output))
					exitCodeChannel <- programExitCode
				},
			)

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, regex, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(result).To(Equal("123"))
		})

		It("should only return a result after the command finishes running", func() {
			executionOrderChannel := make(chan string, 1)

			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(output))
					time.Sleep(100 * time.Millisecond)
					exitCodeChannel <- programExitCode
					executionOrderChannel <- "command finished"
				},
			)

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, regex, expectedProgramArg)
			go func() { executionOrderChannel <- "result returned" }()

			Expect(<-executionOrderChannel).To(Equal("command finished"))
			Expect(err).To(BeNil())
			Expect(result).To(Equal("123"))
		})

		It("should return an empty string if there are no regex matches in the cumulative command output", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(programStdout))
					exitCodeChannel <- programExitCode
				},
			)

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, regex, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(""))
		})

		It("should only consider the stdout command output when finding a regex match", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stderrWriter.Write([]byte(output))
					exitCodeChannel <- programExitCode
				},
			)

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, regex, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(""))
		})

		It("should return an empty string if the regex contains no capturing group", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(output))
					exitCodeChannel <- programExitCode
				},
			)

			regexWithoutGroup, _ := regexp.Compile("The result is \\d+")

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, regexWithoutGroup,
				expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(""))
		})
	})

	When("`resultCaptureRegex` is set to `nil`", func() {
		It("should return an empty string", func() {
			shellCommandDouble = test.NewShellCommandDouble(
				func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
					exitCodeChannel chan<- int) {

					_, _ = stdoutWriter.Write([]byte(programStdout))
					exitCodeChannel <- programExitCode
				},
			)

			result, err := shellCommandService.RunShellCommand(expectedProgram, false, nil, expectedProgramArg)

			Expect(err).To(BeNil())
			Expect(result).To(Equal(""))
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
