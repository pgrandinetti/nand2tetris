package main

import (
	"strings"
)

// Assembly instructions are \n-separated strings, trimmed of leading/trailing spaces.
type Instruction = string

const (
	A_INSTRUCTION = iota
	C_INSTRUCTION
	L_INSTRUCTION
)

func Advance(rawLine string) string {
	// Reads a string from input ('rawLine'), skips whitespaces
	// and returns the string of the current instruction (trimmed).
	// If rawLine is blank or a comment line, returns "".
	//
	// Assume that the assembly code is error-free (A1).
	inpt := strings.TrimSpace(rawLine)
	if len(inpt) == 0 {
		return inpt
	}
	if inpt[0] == '/' {
		// Comment line (A1).
		return ""
	}
	return inpt
}

func InstructionType(instr Instruction) int {
	// Cases based on A1.
	if instr[0] == '@' {
		return A_INSTRUCTION
	}
	if instr[0] == '(' {
		return L_INSTRUCTION
	}
	return C_INSTRUCTION
}

func Symbol(instr Instruction, instrType int) string {
	// If the current instruction is '(xxx)', then return 'xxx'.
	// If the current instruction is '@xxx', then return 'xxx'.
	// Should not be called in other cases.
	if instrType == L_INSTRUCTION {
		return instr[1 : len(instr)-1]
	}
	if instrType == A_INSTRUCTION {
		return instr[1:]
	}
	return "" // Should never happen according to the specifications.
}

func Dest(instr Instruction) string {
	// Should be called only for C_INSTRUCTION.
	// Returns the symbolic 'dest' part of the C-instruction (8 possibilities).
	// dest=comp;jump
	// 'dest' is optional, in which case '=' is omitted.
	for i := 0; i < len(instr); i++ {
		if instr[i] == '=' {
			return instr[0:i]
		}
	}
	return ""
}

func Comp(instr Instruction) string {
	// Should be called only for C_INSTRUCTION.
	// Returns the symbolic 'comp part of the C-instruction (28 possibilities).
	// dest=comp;jump
	// 'comp' is mandatory.
	i := 0
	j := 0
	for ; j < len(instr); j++ {
		if instr[j] == '=' {
			i = j + 1
		}
		if instr[j] == ';' {
			break
		}
	}
	return instr[i:j]
}

func Jump(instr Instruction) string {
	// Should be called only for C_INSTRUCTION.
	// Returns the symbolic 'jump' part of the C-instruction (8 possibilities).
	// dest=comp;jump
	// 'jump' is optional, in which case ';' is omitted.
	for i := 0; i < len(instr); i++ {
		if instr[i] == ';' {
			return instr[i+1:]
		}
	}
	return ""
}
