package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Translate(fpath string) string {
	reader, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()
	basePath := filepath.Base(fpath)
	CurrentVMFile = basePath[:len(basePath)-len(filepath.Ext(basePath))]
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
		b.WriteString("// " + instr + "\n")
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
		//b.WriteString("// [END] " + instr + "\n")
	}
	return b.String()
}

func AddOps() string {
	var b strings.Builder
	//b.WriteString("@END\n0;JMP\n")
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
	var err error
	var program string
	var fout string
	fin := os.Args[1]
	fi, err := os.Stat(fin)
	if err != nil {
		panic(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		filepath.WalkDir(fin, func(s string, d fs.DirEntry, e error) error {
			if e != nil {
				panic(e)
			}
			if filepath.Ext(d.Name()) == ".vm" {
				program += fmt.Sprintf("\n// [FILE] %s\n", d.Name())
				program += Translate(filepath.Join(fin, d.Name()))
			}
			return nil
		})
		fout = fmt.Sprintf("%s.asm", fin)
	case mode.IsRegular():
		program = Translate(fin)
		fout = fmt.Sprintf("%sasm", fin[:len(fin)-2])
	}
	if !strings.Contains(program, "Sys.init") {
		panic("The program does not contain a Sys.init function")
	}
	program = WriteBootstrap() + program + AddOps()
	err = os.WriteFile(fout, []byte(program), 0644)
	if err != nil {
		panic(err)
	}
}
