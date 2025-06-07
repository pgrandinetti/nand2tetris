package main

// symTable holds the state of the Symbol Table for the Assembler.
// The function AddEntry has the side effect to modify the state adding new symbols.
var symTable = map[string]int{
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"SCREEN": 16384,
	"KBD":    24576}

// According to the specifications, when a new symbol is defined it is assigned
// a memory address starting from 16.
// The function NextAvailAddr has the side effect to increment this counter.
var nextAddr int = 16

func NextAvailAddr() int {
	ret := nextAddr
	nextAddr++
	return ret
}

func AddEntry(sym string, address int) {
	symTable[sym] = address
}

func Contains(sym string) bool {
	_, ok := symTable[sym]
	return ok
}

func GetAddress(sym string) int {
	return symTable[sym]
}
