package main

import (
    "fmt"
    "log"
    "os"
    "bufio"
    "strings"
)

func Translate(fpath string) string {
    reader, err := os.Open(fpath)
    if err != nil {
        log.Fatal(err)
    }
    defer reader.Close()
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
        b.WriteString("// " + instr + "\n")
        cType = CommandType(instr)
        if cType == C_ARITHMETIC {
            b.WriteString(WriteArithmetic(instr))
        } else {
            b.WriteString(WritePushPop(cType, Arg1(instr, cType), Arg2(instr)))
        }
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
