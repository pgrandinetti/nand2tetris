package main

import (
	"os"
)

func main() {
	fname := os.Args[1]
	content, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	// This will print on stdout
	CompileClass(content, 0)
}
