package main

import (
    "fmt"
    "strings"
)

// pc is used to make new labels, eg. (LABEL.0), (LABEL.1), ...
var pc = 0

// Push value from D register onto stack (RAM[SP++] = D)
const pushDtoStack = "  //PUSH\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"

// D = RAM[--SP]
const popStackToD = "  //POP\n@SP\nM=M-1\nA=M\nD=M\n"

// All computatins below will jump back to where they were called from
const procDone = "@DONE\n0;JMP\n"

// Operations add, sub, and, or utilize R13 to save the first operand (stack -> D -> R13)
// and then implement D = D `op` R13
const add = "(ADD)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D+M\n" + pushDtoStack + procDone
const sub = "(SUB)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=M-D\n" + pushDtoStack + procDone  // in M-D order matters
const and = "(AND)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D&M\n" + pushDtoStack + procDone
const or = "(OR)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D|M\n" + pushDtoStack + procDone

// Operations neg, not perform
// D = op D
const neg = "(NEG)\n" + popStackToD + "D=-D\n" + pushDtoStack + procDone
const not = "(NOT)\n" + popStackToD + "D=!D\n" + pushDtoStack + procDone

// Operations eq, lt, gt utilize R13 for the first operand,
// then compute D-R13 and execute a JMP instruction accordingly. Then
// D is set to either -1 (true) or 0 (false).
const eq = "(EQ)\n" + popStackToD + "  // save to R13\n@R13\nM=D\n" + popStackToD + `  // compute
@R13
D=D-M
@eq.TRUE
D;JEQ
(eq.FALSE)
D=0
@eq.DONE
0;JMP
(eq.TRUE)
D=-1
(eq.DONE)
` + pushDtoStack + procDone

const lt = "(LT)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + `
@R13
D=D-M
@lt.TRUE
D;JLT
(lt.FALSE)
D=0
@lt.DONE
0;JMP
(lt.TRUE)
D=-1
(lt.DONE)
` + pushDtoStack + procDone

const gt = "(GT)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + `
@R13
D=D-M
@gt.TRUE
D;JGT
(gt.FALSE)
D=0
@gt.DONE
0;JMP
(gt.TRUE)
D=-1
(gt.DONE)
` + pushDtoStack + procDone

// Used to jump back after executing a add,neg,gt,... computation.
// This uses a label "BACK" that needs to be set as a pointer to the real
// "back" instruction (BACK.0, BACK.1, ...)
const done = "(DONE)\n@BACK\nA=M\n0;JMP\n"

// end implements an infinite loop
const end = "(END)\n@END\n0;JMP\n"

// D = segment[index]
func segmentToD(seg string, index int) string {
    if seg == "local" {
        return fmt.Sprintf("@%d\nD=A\n@LCL\nA=D+M\nD=M\n", index)
    }
    if seg == "argument" {
        return fmt.Sprintf("@%d\nD=A\n@ARG\nA=D+M\nD=M\n", index)
    }
    if seg == "this" {
        return fmt.Sprintf("@%d\nD=A\n@THIS\nA=D+M\nD=M\n", index)
    }
    if seg == "that" {
        return fmt.Sprintf("@%d\nD=A\n@THAT\nA=D+M\nD=M\n", index)
    }
    if seg == "pointer" {
        if index == 0 {
            return "@THIS\nD=M\n"
        }
        // else must be index=1
        return "@THAT\nD=M\n"
    }
    if seg == "temp" {
        // index must be [0, 7]
        return fmt.Sprintf("@5\nA=A+%d\nD=M\n", index)
    }
    if seg == "constant" {
        return fmt.Sprintf("@%d\nD=A\n", index)
    }
    // Must be "static"
    // The textbook suggests to use the variable 'FileName.index'
    // but I return a fixed 'Prog.index'.
    return fmt.Sprintf("@Prog.%d\nA=M\nD=M\n", index)
}

func WriteArithmetic(instr Instruction) string {
    /* Make a new label (BACK.pc),
       then set the value of (BACK) to this new label because it's used in (DONE)
       then jump to add/sub/neg/...
    Eg. if pc=53
     @BACK.53
     D=A
     @BACK
     M=D      // now (BACK) points to (BACK.53) and will be jumped to
              // from (DONE) which is a goto for all ADD/SUB/NEG operations
     @ADD
     0;JMP
     (BACK.53)
    */
    var b strings.Builder
    op := Arg1(instr, C_ARITHMETIC)
    label := fmt.Sprintf("BACK.%d", pc)
    b.WriteString(fmt.Sprintf("@%s\n", label))
    b.WriteString("D=A\n")
    b.WriteString("@BACK\n")
    b.WriteString("M=D\n")
    b.WriteString(fmt.Sprintf("@%s\n", strings.ToUpper(op)))
    b.WriteString("0;JMP\n")
    b.WriteString(fmt.Sprintf("(%s)\n", label))
    pc++
    return b.String()
}

func WritePushPop(cType int, seg string, index int) string {
    var b strings.Builder
    b.WriteString(segmentToD(seg, index))
    if cType == C_PUSH {
        b.WriteString(pushDtoStack)
    } else {
        b.WriteString(popStackToD)
    }
    return b.String()
}

func main() {
    lines := [...]string {"push constant 7", "push local 2", "add", "push constant 13", "sub", "push constant 3", "eq"}
    var b strings.Builder
    var cType int
    var instr Instruction
    for _, elem := range lines {
        b.WriteString("// " + elem + "\n")
        instr = Advance(elem)
        cType = CommandType(instr)
        if cType == C_ARITHMETIC {
            b.WriteString(WriteArithmetic(instr))
        } else {
            b.WriteString(WritePushPop(cType, Arg1(instr, cType), Arg2(instr)))
        }
    }
    b.WriteString("@END\n0;JMP\n")
    b.WriteString(add)
    b.WriteString(sub)
    b.WriteString(eq)
    b.WriteString(done)
    b.WriteString(end)
    fmt.Println(b.String())
}
