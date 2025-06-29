package main

import (
	"strconv"
	"strings"
)

// VM instructions are \n-separated strings, trimmed of leading/trailing spaces
// Empty instructions ("") correspond to blank or comment lines in the source code.
type Instruction = string

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
	UNDEFINED
)

func Advance(rawLine string) Instruction {
	// Reads a string from input ('rawLine'), skips whitespaces
	// and returns the string of the current instruction (trimmed).
	// If rawLine is blank or a comment line, returns "".
	//
	// Assume that the VM code is error-free (A1).

	inpt := strings.TrimSpace(rawLine)
	if len(inpt) == 0 || inpt[0:2] == "//" {
		// Empty or comment line (A1).
		return ""
	}

	i := 0 // index of first '/'
	for ; i < len(inpt); i++ {
		if inpt[i] == '/' {
			break
		}
	}
	inpt = strings.TrimSpace(inpt[:i])
	if len(inpt) == 0 {
		return inpt
	}
	return strings.Join(strings.Fields(inpt), " ") // collapse whitespaces into one
}

func CommandType(instr Instruction) int {
	if len(instr) > 4 && instr[0:4] == "push" {
		return C_PUSH
	}
	if len(instr) > 3 && instr[0:3] == "pop" {
		return C_POP
	}
	if instr == "add" || instr == "sub" || instr == "neg" || instr == "eq" || instr == "gt" || instr == "lt" || instr == "and" || instr == "or" || instr == "not" {
		return C_ARITHMETIC
	}
	if len(instr) > 5 && instr[0:5] == "label" {
		return C_LABEL
	}
	if len(instr) > 4 && instr[0:4] == "goto" {
		return C_GOTO
	}
	if len(instr) > 7 && instr[0:7] == "if-goto" {
		return C_IF
	}
	if len(instr) > 8 && instr[0:8] == "function" {
		return C_FUNCTION
	}
	if instr == "return" {
		return C_RETURN
	}
	if len(instr) > 4 && instr[0:4] == "call" {
		return C_CALL
	}
	return UNDEFINED
}

func Arg1(instr Instruction, cType int) string {
	if cType == C_ARITHMETIC {
		return instr
	}
	s1 := 0 // There will always be at least one blank space
	s2 := 0 // There may not be a second blank space
	for ; s2 < len(instr); s2++ {
		if instr[s2] == ' ' {
			if s1 == 0 {
				s1 = s2
			} else {
				break
			}
		}
	}
	return instr[s1+1 : s2]
}

func Arg2(instr Instruction) int {
	j := 0
	for j = len(instr) - 1; j >= 0; j-- {
		if instr[j] == ' ' {
			break
		}
	}
	i, _ := strconv.Atoi(instr[j+1:])
	return i
}
