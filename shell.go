package gsh

import (
    "errors"
    "fmt"
    "os"
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
        fmt.Println(currentCommand)
        success := false
        builtin := shell.getBuiltIn(currentCommand.executable)
        if builtin != nil {
            success = shell.execInternal(builtin)
        } else {
            success = execCommand(currentCommand)
        }

        if currentCommand.next != nil {
             if canContinue(success, currentCommand.nextCombinator) {
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

func canContinue(processSuccess bool, combinator combinatorType) bool {
    if combinator == COMBINATOR_OR {
        return processSuccess == false
    } else if combinator == COMBINATOR_AND {
        return processSuccess == true
    }

    return true
}

func execCommand(command *command) bool {
    execCmd := command.ToExecCommand()
    err := execCmd.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error(), "\n")
    }

    return execCmd.ProcessState.Success()
}

func (shell *Shell) execInternal(builtin Builtin) bool {
    builtin.Execute(shell)
    return true
}
