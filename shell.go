package gsh

import (
    "errors"
    "fmt"
    "os"
    "os/exec"
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

func (shell *Shell) getBuiltIn(name string) Builtin {
    for _, builtin := range shell.builtins {
        if (name == builtin.GetName()) {
            return builtin
        }
    }

    return nil
}

func (shell *Shell) execCommandChain(command *command) {
    currentCommand := command

    for {
        execCmd := currentCommand.ToExecCommand()
        err := execCmd.Run()
        if err != nil {
            fmt.Fprintf(os.Stderr, err.Error(), "\n")
        }

        if currentCommand.next != nil {
             if canContinue(execCmd, currentCommand.nextCombinator) {
                 currentCommand = currentCommand.next
                 continue;
             } else {
                 break;
             }
        } else {
            break;
        }
    }
}

func canContinue(execCmd *exec.Cmd, combinator combinatorType) bool {
    processSuccess := execCmd.ProcessState.Success()

    if combinator == COMBINATOR_OR {
        return processSuccess == false
    } else if combinator == COMBINATOR_AND {
        return processSuccess == true
    }

    return true
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
