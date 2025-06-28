package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
    "path/filepath"
)

func Translate(fpath string) string {
	reader, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
    CurrentVMFile = filepath.Base(fpath)
	var b strings.Builder
	scanner := bufio.NewScanner(reader)
	var line string
	var cType int
	var instr Instruction
	for scanner.Scan() {
		line = scanner.Text()
		instr = Advance(line)
		if len(instr) == 0 {
			continue
		}
		cType = CommandType(instr)
        b.WriteString("// " + instr + fmt.Sprintf("(type %d)\n", cType))
		if cType == C_ARITHMETIC {
			b.WriteString(WriteArithmetic(instr))
        } else if cType == C_LABEL {
            b.WriteString(WriteLabel(Arg1(instr, cType)))
        } else if cType == C_GOTO {
            b.WriteString(WriteGoto(Arg1(instr, cType)))
        } else if cType == C_IF {
            b.WriteString(WriteIf(Arg1(instr, cType)))
        } else if cType == C_FUNCTION {
            CurrentVMFunc = Arg1(instr, cType)
            b.WriteString(WriteFunction(Arg1(instr, cType), Arg2(instr)))
        } else if cType == C_CALL {
            b.WriteString(WriteCall(Arg1(instr, cType), Arg2(instr)))
        } else if cType == C_RETURN {
            b.WriteString(WriteReturn())
        } else if cType == C_PUSH || cType == C_POP {
			b.WriteString(WritePushPop(cType, Arg1(instr, cType), Arg2(instr)))
		} else {
            msg := fmt.Sprintf("cType not valid: %s", instr) 
            panic(msg)
        }
        b.WriteString("// [END] " + instr + "\n")
	}
	b.WriteString("@END\n0;JMP\n")
	b.WriteString("//\n// PROCEDURES SECTION\n//\n")
	b.WriteString(add)
	b.WriteString(sub)
	b.WriteString(neg)
	b.WriteString(gt)
	b.WriteString(lt)
	b.WriteString(eq)
	b.WriteString(and)
	b.WriteString(or)
	b.WriteString(not)
	b.WriteString(done)
	b.WriteString(end)
	return b.String()
}

func main() {
	fin := os.Args[1]
	program := Translate(fin)
	fout := fmt.Sprintf("%sasm", fin[:len(fin)-2])
	err := os.WriteFile(fout, []byte(program), 0644)
	if err != nil {
		panic(err)
	}
}
