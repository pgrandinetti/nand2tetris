package main

import (
	"fmt"
	"os"
)

// Run this on .jack files one-by-one. Eg.:
// go build -o a.out
// ./a.out Pong/PongGame.jack > Pong/PongGame.vm
// ./a.out ... # all files in the 'Pong' directory
// then load the directory 'Pong' in the VM emulator.
func main() {
	fname := os.Args[1]
	content, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	_, code := CompileClass(content, 0)
	fmt.Println(code)
}
