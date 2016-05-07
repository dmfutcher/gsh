package gsh

type Builtin interface {
    GetName()       string
    Execute(*Shell) error
}

type BuiltinExit struct {}

func (_ BuiltinExit) GetName() string {
    return "exit"
}

func (_ BuiltinExit) Execute(shell *Shell) error {
    shell.exiting = true
    return nil
}
