package gsh

import (
    "fmt"
    "os"
    "os/exec"
)

type combinatorType int
const (
    COMBINATOR_NULL = iota
    COMBINATOR_OR
    COMBINATOR_AND
    COMBINATOR_UNCONDITIONAL
)

type output string
type input string
const (
    OUTPUT_STDOUT = "stdout"
    OUTPUT_STDERR = "stderr"
    INPUT_STDIN = "stdin"
)

type command struct {
    executable      string
    args            []string

    input           input
    output          output
    errOutput       output

    next            *command
    nextCombinator  combinatorType
}

func newCommand(executable string, args []string) *command {
    return &command{
        executable: executable,
        args: args,
        input: INPUT_STDIN,
        output: OUTPUT_STDOUT,
        errOutput: OUTPUT_STDOUT,
    }
}

func (command command) ToExecCommand() *exec.Cmd {
    execCmd := exec.Command(command.executable, command.args...)

    output := getOutputHandle(command.output)
    if output != nil {
        execCmd.Stdout = output
    }

    errOutput := getOutputHandle(command.errOutput)
    if errOutput != nil {
        execCmd.Stderr = errOutput
    }

    input := getInputHandle(command.input)
    if input != nil {
        execCmd.Stdin = input
    }

    return execCmd
}

func getOutputHandle(output output) *os.File {
    switch (output) {
    case OUTPUT_STDOUT:
        return os.Stdout
    case OUTPUT_STDERR:
        return os.Stderr
    default:
        fmt.Fprintf(os.Stderr, "Non-standard outputs not supported")
        return nil
    }
}

func getInputHandle(input input) *os.File {
    switch (input) {
    case INPUT_STDIN:
        return os.Stdin
    default:
        fmt.Fprintf(os.Stderr, "Non-standard outputs not supported")
        return nil
    }
}
