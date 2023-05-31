package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Exec(bin string, args ...string) (result string, err error) {
	cmd := exec.Command(bin, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

func ExecQuit(bin string, args ...string) (result string, err error) {
	cmd := exec.Command(bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return "", err
	}
	fmt.Printf(bin, args)
	err = cmd.Wait()
	return out.String(), err
}
