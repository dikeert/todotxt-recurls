package main

import (
	"errors"
	"os"
	"os/exec"
)

func CreateCmd() (*exec.Cmd, error) {
	if len(os.Args) > 1 {
		cmd := exec.Command("todo.sh", "command")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		for _, arg := range os.Args[1:] {
			cmd.Args = append(cmd.Args, arg)
		}

		return cmd, nil
	} else {
		return nil, errors.New("command is missing")
	}
}
