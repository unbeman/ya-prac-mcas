package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello, World!")
	os.Exit(1) // want "call of os.Exit in main function"
}
