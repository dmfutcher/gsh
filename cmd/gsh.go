package main

import (
    "fmt"
    "os"

    "github.com/bobbo/gsh"
)

func main() {
    config := gsh.ShellConfig{Interactive: true}
    shell := gsh.New(config)

    code, err := shell.Run()
    if err != nil {
        fmt.Fprintf(os.Stderr, err.Error())
    }

    os.Exit(code)
}
