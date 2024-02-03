package main

import (
	"fmt"
	"miniredis/src/server"
)

func main() {
	fmt.Println("Hello, Modules!")

	err := server.StartServer("127.0.0.1", 5093)

	if err != nil {
		fmt.Println(err)
	}

	// Parser := common.NewParser()

	// Parser.Parse([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))
}
