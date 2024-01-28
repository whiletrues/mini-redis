package server

import (
	"fmt"
	"miniredis/src/common"
	"net"
)

type server struct {
	listener net.Listener
}

func StartServer() {
	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	for {

		con, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(con)
	}
}

func handleConnection(con net.Conn) {
	defer con.Close()

	buffer := make([]byte, 1024)

	_, err := con.Read(buffer)

	if err != nil {
		fmt.Println(err)
		return
	}

	common.NewParser()

}
