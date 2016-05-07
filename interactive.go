package gsh

import (
    "bufio"
    "fmt"
    "os"
)

type interactiveShell struct {
    Shell
    prompt  string
}

func interactiveMainLoop(shell Shell) (int, error) {
    interactive := interactiveShell{Shell: shell, prompt: "> "}
    return interactive.mainLoop()
}

func (shell *interactiveShell) mainLoop() (int, error) {
    for !shell.exiting {
        input, _ := shell.readInput()
        if input == "" {
            continue
        }        
        command := parseInput(input)
        shell.execCommand(command)
    }

    return 0, nil
}

func (shell *interactiveShell) readInput() (string, error) {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print(shell.prompt)
    scanner.Scan()
    return scanner.Text(), nil
}
