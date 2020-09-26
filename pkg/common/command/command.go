package command

import (
	"io"
	"os/exec"
)

type Command struct {
	*exec.Cmd
}

func NewCommand(cmd *exec.Cmd, stdin io.Reader, stdout io.Writer, stderr io.Writer) (*Command, error) {
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Command{Cmd: cmd}, nil
}