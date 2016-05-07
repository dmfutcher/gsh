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
    config  *ShellConfig
    exiting bool
}

func ShellRun(config ShellConfig) (int, error) {
    shell := Shell{&config, false}

    if !config.Interactive {
        return -1, errors.New("Non-interactive not implemented")
    } else {
        return interactiveMainLoop(shell)
    }
}

func parseInput(input string) *exec.Cmd {
    tokens := strings.Split(input, " ")
    return exec.Command(tokens[0], tokens[1:]...)
}

func (shell Shell) execCommand(cmd *exec.Cmd) {
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err.Error())
    }
}
