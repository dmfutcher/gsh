package gsh

import (
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

type ShellConfig struct {
    Interactive bool
}

type Shell struct {
    config      *ShellConfig
    exiting     bool
    builtins    []Builtin
}

func New(config ShellConfig) Shell {
    shell := Shell{
        config: &config,
        exiting: false,
        builtins: []Builtin{
            BuiltinExit{},
        },
    }

    return shell
}

func (shell *Shell) Run() (int, error) {
    if !shell.config.Interactive {
        return -1, errors.New("Non-interactive not implemented")
    } else {
        return interactiveMainLoop(shell)
    }
}

func parseInput(input string) *exec.Cmd {
    tokens := strings.Split(input, " ")
    return exec.Command(tokens[0], tokens[1:]...)
}

func (shell *Shell) getBuiltIn(name string) Builtin {
    for _, builtin := range shell.builtins {
        if (name == builtin.GetName()) {
            return builtin
        }
    }

    return nil
}

func (shell *Shell) execCommand(cmd *exec.Cmd) {
    builtin := shell.getBuiltIn(cmd.Path)
    if builtin != nil {
        builtin.Execute(shell)
        return
    }

    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    }
}
