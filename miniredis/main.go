package main

import (
	"fmt"
	netpoller "miniredis/src/server"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("Hello, Modules!")

	netpoller, err := netpoller.New()

	if err != nil {
		fmt.Println(err)
		return
	}

	r, w, err := socketPair()
	
	netpoller.Start()
}
func socketPair() (r, w int, err error) {
	fd, err := unix.Socketpair(unix.AF_UNIX, unix.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	if err = unix.SetNonblock(fd[0], true); err != nil {
		return
	}
	if err = unix.SetNonblock(fd[1], true); err != nil {
		return
	}

	buf := 4096
	if err = unix.SetsockoptInt(fd[0], unix.SOL_SOCKET, unix.SO_SNDBUF, buf); err != nil {
		return
	}
	if err = unix.SetsockoptInt(fd[1], unix.SOL_SOCKET, unix.SO_SNDBUF, buf); err != nil {
		return
	}
	if err = unix.SetsockoptInt(fd[0], unix.SOL_SOCKET, unix.SO_RCVBUF, buf); err != nil {
		return
	}
	if err = unix.SetsockoptInt(fd[1], unix.SOL_SOCKET, unix.SO_RCVBUF, buf); err != nil {
		return
	}
	return fd[0], fd[1], nil
}
