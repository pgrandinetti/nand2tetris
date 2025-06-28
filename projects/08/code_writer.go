package main

import (
	"fmt"
	"strings"
)

// The name of the .vm file that is currently being translated is needed
// to make labels and static variables names.
var CurrentVMFile = ""

// The name of the 'function' body that is currently being translated is used to make
// labels unique.
var CurrentVMFunc = ""

// pc is used (and incremented) to make new labels, eg. (LABEL.0), (LABEL.1), ...
var pc = 0

var segNames = map[string]string{"local": "LCL", "argument": "ARG", "this": "THIS", "that": "THAT", "temp": "TEMP", "constant": "constant"}

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
    //seg = segNames[seg]
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
    if seg == "static" {
	    return fmt.Sprintf("@%s.%d\nD=M\n", CurrentVMFile, index)
    }
    msg := fmt.Sprintf("not valid: segment |%s|", seg)
    panic(msg)
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
    if seg == "static" {
	    // else must be "static", because "constant" is virtual.
	    return fmt.Sprintf("@%s.%d\nM=D\n", CurrentVMFile, index)
    }
    msg := fmt.Sprintf("segment not valid: |%s|", seg)
    panic(msg)
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

func WriteLabel(arg1 string) string {
	return fmt.Sprintf("  // label\n(%s.%s$%s)\n", CurrentVMFile, CurrentVMFunc, arg1)
}

func WriteGoto(arg1 string) string {
	return fmt.Sprintf("  // goto\n@%s.%s$%s\n0;JMP\n", CurrentVMFile, CurrentVMFunc, arg1)
}

func WriteIf(arg1 string) string {
	var b strings.Builder
	b.WriteString("  // if-goto\n")
	b.WriteString(popStackToD)
	b.WriteString(fmt.Sprintf("@%s.%s$%s\n", CurrentVMFile, CurrentVMFunc, arg1))
	b.WriteString("D;JNE\n")
	return b.String()
}

func WriteFunction(fName string, nVars int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  // function %s %d\n", fName, nVars))
	b.WriteString(fmt.Sprintf("(%s.%s)\n", CurrentVMFile, fName))
	for range nVars {
		// Initialize local variables for the function.
		// repeat 'nVars' times:
		// push constant 0
		b.WriteString(WritePushPop(C_PUSH, "constant", 0))
	}
	return b.String()
}

func WriteCall(fName string, nArgs int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("  // call %s %d\n", fName, nArgs))
	retAddr := fmt.Sprintf("%s.%s$ret.%d", CurrentVMFile, CurrentVMFunc, pc)
	pc++
	// push returnAddress
	b.WriteString(fmt.Sprintf("@%s\nD=A\n", retAddr) + pushDtoStack)
	// push LCL
	b.WriteString("@LCL\nD=M\n" + pushDtoStack)
	// push ARG
	b.WriteString("@ARG\nD=M\n" + pushDtoStack)
	// push THIS
	b.WriteString("@THIS\nD=M\n" + pushDtoStack)
	// push THAT
	b.WriteString("@THAT\nD=M\n" + pushDtoStack)
	// ARG = SP - 5 - nArgs
	b.WriteString(fmt.Sprintf("@5\nD=A\n@%d\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n", nArgs))
	// LCL = SP
	b.WriteString("@SP\nD=M\n@LCL\nM=D\n")
	// goto fName
	b.WriteString(fmt.Sprintf("@%s.%s\n0;JMP\n", CurrentVMFile, fName))
	// insert retAddr now
	b.WriteString("(" + retAddr + ")\n")
	return b.String()
}

func WriteReturn() string {
	var b strings.Builder
	b.WriteString("  // return\n")
	// R13 = LCL (use R13 to save the address of the function frame)
	b.WriteString("@LCL\nD=M\n@R13\nM=D\n")
	// retAddr = *(R13 - 5) (use R14 to store this)
	b.WriteString("@5\nD=A\n@R13\nD=M-D\nA=D\nD=M\n@R14\nM=D\n")
	// *ARG = pop()
	b.WriteString(popStackToD + "@ARG\nA=M\nM=D\n")
	// SP = ARG + 1
	b.WriteString("@ARG\nD=M+1\n@SP\nM=D\n")
	// THAT = *(R13 - 1)
	b.WriteString("@R13\nD=M-1\nA=D\nD=M\n@THAT\nM=D\n")
	// THIS = *(R13 - 2)
	b.WriteString("@2\nD=A\n@R13\nD=M-D\nA=D\nD=M\n@THIS\nM=D\n")
	// ARG = *(R13 - 3)
	b.WriteString("@3\nD=A\n@R13\nD=M-D\nA=D\nD=M\n@ARG\nM=D\n")
	// LCL = *(R13 - 4)
	b.WriteString("@4\nD=A\n@R13\nD=M-D\nA=D\nD=M\n@LCL\nM=D\n")
	// goto retAddr
	b.WriteString("@R14\nA=M\n0;JMP\n")
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
