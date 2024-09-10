package netkit

import (
	"net"
)

func RandomPort() (port int, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, err
}
