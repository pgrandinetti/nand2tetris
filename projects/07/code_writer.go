package main

import (
	"fmt"
	"strings"
)

// pc is used (and incremented) to make new labels, eg. (LABEL.0), (LABEL.1), ...
var pc = 0

var segNames = map[string]string{"local": "LCL", "argument": "ARG", "this": "THIS", "that": "THAT", "temp": "TEMP"}

// RAM[SP++] = D
const pushDtoStack = "  //PUSH\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"

// D = RAM[--SP]
const popStackToD = "  //POP\n@SP\nM=M-1\nA=M\nD=M\n"

/*
   The arithmetic-logical operation (below) are implemented so that they are independent from
   where in the RAM they are executed. For instance, 'add' does the following:
     1. pop from stack to D reg. Save D to R13.
     2. pop from stack to D.
     3. compute D + R13. Push result to stack.
     4. jump to (DONE).
   All these operations are implemented so to jump to (DONE) eventually. The point being to avoid duplicating code.
   This is further explained below.

   Consider the VM program
   ...
   add
   add

   If translated line-by-line, the two lines 'add' would generate the same code, repeated twice.
   Instead, I use the following ideas:
     - whenever a arithmetic-logical operation is found, a new label is created as 'BACK.pc++' (pc is an integer 0,...)
       For example, the first time a 'add' instruction is encountered, a label 'BACK.0' will be made.
     - the code sets the value of the label 'BACK' to the address of the label 'BACK.0', then it jumps to the label (ADD)
       which is unique in the program.
     - ADD (like all other operations), will do a jump (DONE)
     - DONE loads in the A reg the value pointed to by BACK, which was set to BACK.0. And that's how the PC goes back to where it was.

   Example:

             |  @BACK.0   // create new label, will be any address starting from 16,...
   add ------+  D=A       // save the address of BACK.0 in D
             |  @BACK     // put into A the address of BACK
             |  M=D       // put in the address pointed by BACK the value of D (which now is the value of BACK.0)
             |  @ADD      // jump to the ADD code
             |  0;JMP
             |  (BACK.0)  // ensure the PC jumps back here after ADD is finished
   ...
   ...
   ...       |  @BACK.1   // the same as the first add, but using a new label 'BACK.1'
   add ------+  D=A
             |  @BACK
             |  M=D
             |  @ADD
             |  0;JMP
             |  (BACK.1)

   // code for 'add'. This is written only one time.
   (ADD)
   // ...pop, pop, sum, push
   @DONE
   0;JMP

   // code for 'gt'. This is written only one time.
   (GT)
   // ... pop, pop, diff, set D=-1 or D=0
   @DONE
   0;JMP

   // code for 'done'
   (DONE)
   @BACK   // put in A the value of the label BACK. This is never changed by the program, but the content of that register is changed.
   A=M     // put in A the value of the register pointed by BACK. This would be the value BACK.0, then the value BACK.1, etc.
   0;JMP
*/

const procDone = "@DONE\n0;JMP\n"

// NOTE: In 'push x, push y, sub' the sematic is: y<-pop, x<-pop, D<-x-y (order matters)
const add = "(ADD)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D+M\n" + pushDtoStack + procDone
const sub = "(SUB)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D-M\n" + pushDtoStack + procDone
const and = "(AND)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D&M\n" + pushDtoStack + procDone
const or = "(OR)\n" + popStackToD + "@R13\nM=D\n" + popStackToD + "@R13\nD=D|M\n" + pushDtoStack + procDone

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
		return fmt.Sprintf("@5\nD=A\n@%d\nA=D+A\nD=M\n", index)
	}
	if seg == "constant" {
		return fmt.Sprintf("@%d\nD=A\n", index)
	}
	// Must be "static"
	// The textbook suggests to use the variable 'FileName.index'
	// but I return a fixed 'Prog.index'.
	return fmt.Sprintf("@Prog.%d\nD=M\n", index)
}

// segment[index] = D
func DtoSegment(seg string, index int) string {
	// NOTE. Some of these use both R13 and R14 as intermediate variables.
	// R13 is used to hold the value when D will b overwritten by other computations.
	// R14 is used to compute the final address: for example 'local 4' requires to do '@LCL + 4'
	if seg == "local" || seg == "argument" || seg == "this" || seg == "that" {
		return fmt.Sprintf("  // D->R13\n@R13\nM=D\n  // addr->R14\n@%d\nD=A\n@%s\nD=D+M\n@R14\nM=D\n  // R13->addrOfR14\n@R13\nD=M\n@R14\nA=M\nM=D\n", index, segNames[seg])
	}
	if seg == "pointer" {
		if index == 0 {
			return "@THIS\nM=D\n"
		}
		return "@THAT\nM=D\n"
	}
	if seg == "temp" {
		return fmt.Sprintf("  // D->R13\n@R13\nM=D\n  // addr->R14\n@5\nD=A\n@%d\nD=D+A\n@R14\nM=D\n  // R13->addrOfR14\n@R13\nD=M\n@R14\nA=M\nM=D\n", index)
	}
	// else must be "static", because "constant" is virtual.
	return fmt.Sprintf("@Prog.%d\nM=D\n", index)
}

func WriteArithmetic(instr Instruction) string {
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
	if cType == C_PUSH {
		b.WriteString(segmentToD(seg, index))
		b.WriteString(pushDtoStack)
	} else {
		b.WriteString(popStackToD)
		b.WriteString(DtoSegment(seg, index))
	}
	return b.String()
}

/*
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
*/
