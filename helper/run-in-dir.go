package helper

import (
	"os"
	"os/exec"
)

func RunInDir(dir string, cmdName string, args ...string) error {
	c := exec.Command(cmdName, args...)
	c.Dir = dir
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	return c.Run()
}
