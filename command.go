package gsh

type combinatorType int
const (
    COMBINATOR_NULL = iota
    COMBINATOR_OR
    COMBINATOR_AND
)

const (
    OUTPUT_STDOUT = "stdout"
    OUTPUT_STDERR = "stderr"
    INPUT_STDIN = "stdin"
)

type command struct {
    executable      string
    args            []string

    input           string
    output          string
    errOutput       string

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
