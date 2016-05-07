package main

import (
    "fmt"
    "os"

    "github.com/bobbo/gsh"
)

func main() {
    config := gsh.ShellConfig{Interactive: true}

    code, err := gsh.ShellRun(config)
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
    }

    os.Exit(code)
}
