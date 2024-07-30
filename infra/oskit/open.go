package oskit

import (
	"fmt"
	"os/exec"
	"runtime"
	"syscall"
)

var commands = map[string]string{
	"windows": "cmd.exe",
	"darwin":  "open",
	"linux":   "xdg-open",
}

var Version = "0.1.0"

// Open calls the OS default program for uri
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}
	if runtime.GOOS == "windows" {
		uri = " /c start " + uri
	}
	cmd := exec.Command(run, uri)
	//fmt.Println(run, uri)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
