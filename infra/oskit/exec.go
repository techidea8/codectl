package oskit

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
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

func ExecWithChanel(ctx context.Context, bin string, args []string) (resultch, errorch, stopch chan string) {
	stopch = make(chan string, 1)
	resultch = make(chan string, 10)
	errorch = make(chan string, 10)
	cmd := exec.Command(bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stopch <- err.Error()
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		stopch <- err.Error()
		return
	}
	err = cmd.Start()
	if err != nil {
		stopch <- err.Error()
		return
	}
	readerstd := bufio.NewReader(stdout)
	readererr := bufio.NewReader(stderr)

	//
	go func() {
		for {
			select {
			case <-stopch:
				return
			default:
				line, err2 := readerstd.ReadString('\n')
				if err2 != nil {
					if io.EOF == err2 {
						stopch <- "quit"
					} else {
						stopch <- err.Error()
					}
				}
				resultch <- line
			}

		}

	}()

	go func() {
		for {
			select {
			case <-stopch:
				return
			default:
				line, err2 := readererr.ReadString('\n')
				if err2 != nil || io.EOF == err2 {
					return
				}
				errorch <- line
			}
		}
	}()
	go cmd.Wait()
	return
}
