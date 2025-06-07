package main

import (
	"fmt"
	"os"
)

func main() {
	fin := os.Args[1]
	fout := fmt.Sprintf("%shack", fin[:len(fin)-3])
	FirstPass(fin)
	binary := SecondPass(fin)
	err := os.WriteFile(fout, []byte(binary), 0644)
	if err != nil {
		panic(err)
	}
}
