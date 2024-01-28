package main

import (
	"fmt"
	"miniredis/src/server"
)

func main() {
	server.StartServer()
	fmt.Println("Hello, Modules!")
}
