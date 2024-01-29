package main

import (
	"fmt"
	"miniredis/src/common"
)

func main() {
	fmt.Println("Hello, Modules!")
	Parser := common.NewParser()

	Parser.Parse([]byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"))

	
}
