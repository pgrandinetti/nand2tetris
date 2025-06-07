package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func FirstPass(fpath string) {
	reader, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	nline := 0
	var line string
	var sym string
	var iType int
	for scanner.Scan() {
		line = scanner.Text()
		line = Advance(line)
		if len(line) == 0 {
			continue
		}
		iType = InstructionType(line)
		if iType == L_INSTRUCTION {
			sym = Symbol(line, iType)
			AddEntry(sym, nline)
		} else {
			nline++
		}
	}
}

func SecondPass(fpath string) string {
	reader, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	var line string
	var binLine string
	var iType int
	var output strings.Builder
	for scanner.Scan() {
		line = scanner.Text()
		line = Advance(line)
		if len(line) == 0 {
			continue
		}
		iType = InstructionType(line)
		if iType == L_INSTRUCTION {
			continue
		}
		if iType == C_INSTRUCTION {
			binLine = instrC(line)
		} else {
			// A_INSTRUCTION
			binLine = instrA(line)
		}
		output.WriteString(binLine)
		output.WriteByte('\n')
	}
	return output.String()
}

func instrC(line string) string {
	dest := Dest(line)
	comp := Comp(line)
	jump := Jump(line)
	return fmt.Sprintf("111%s%s%s", CompCode(comp), DestCode(dest), JumpCode(jump))
}

func instrA(line string) string {
	sym := Symbol(line, A_INSTRUCTION)
	symValue, err := strconv.Atoi(sym)
	if err != nil {
		// Is not integer
		found := Contains(sym)
		if found {
			symValue = GetAddress(sym)
		} else {
			symValue = NextAvailAddr()
			AddEntry(sym, symValue)
		}
	}
	return fmt.Sprintf("0%s", fmt.Sprintf("%015b", symValue)[0:15])
}
