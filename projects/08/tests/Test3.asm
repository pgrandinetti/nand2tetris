// push constant 2(type 1)
@2
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 2
// push constant 3(type 1)
@3
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 3
// call mult 2(type 8)
  // call mult 2
@Test3.vm.$ret.0
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@LCL
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@ARG
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@THIS
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@THAT
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@5
D=A
@2
D=D+A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@Test3.vm.mult
0;JMP
(Test3.vm.$ret.0)
// [END] call mult 2
// push constant 5(type 1)
@5
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 5
// add(type 0)
@BACK.1
D=A
@BACK
M=D
@ADD
0;JMP
(BACK.1)
// [END] add
// pop static 0(type 2)
  //POP
@SP
M=M-1
A=M
D=M
@Test3.vm.0
M=D
// [END] pop static 0
// label END(type 3)
  // label
(Test3.vm.$END)
// [END] label END
// goto END(type 4)
  // goto
@Test3.vm.$END
0;JMP
// [END] goto END
// function mult 2(type 6)
  // function mult 2
(Test3.vm.mult)
@0
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] function mult 2
// push constant 0(type 1)
@0
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 0
// pop local 0(type 2)
  //POP
@SP
M=M-1
A=M
D=M
  // D->R13
@R13
M=D
  // addr->R14
@0
D=A
@LCL
D=D+M
@R14
M=D
  // R13->addrOfR14
@R13
D=M
@R14
A=M
M=D
// [END] pop local 0
// push constant 0(type 1)
@0
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 0
// pop local 1(type 2)
  //POP
@SP
M=M-1
A=M
D=M
  // D->R13
@R13
M=D
  // addr->R14
@1
D=A
@LCL
D=D+M
@R14
M=D
  // R13->addrOfR14
@R13
D=M
@R14
A=M
M=D
// [END] pop local 1
// label WHILE_LOOP(type 3)
  // label
(Test3.vm.mult$WHILE_LOOP)
// [END] label WHILE_LOOP
// push local 1(type 1)
@1
D=A
@LCL
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push local 1
// push argument 1(type 1)
@1
D=A
@ARG
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push argument 1
// lt(type 0)
@BACK.2
D=A
@BACK
M=D
@LT
0;JMP
(BACK.2)
// [END] lt
// not(type 0)
@BACK.3
D=A
@BACK
M=D
@NOT
0;JMP
(BACK.3)
// [END] not
// if-goto WHILE_END(type 5)
  // if-goto
  //POP
@SP
M=M-1
A=M
D=M
@Test3.vm.mult$WHILE_END
D;JNE
// [END] if-goto WHILE_END
// push local 0(type 1)
@0
D=A
@LCL
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push local 0
// push argument 0(type 1)
@0
D=A
@ARG
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push argument 0
// add(type 0)
@BACK.4
D=A
@BACK
M=D
@ADD
0;JMP
(BACK.4)
// [END] add
// pop local 0(type 2)
  //POP
@SP
M=M-1
A=M
D=M
  // D->R13
@R13
M=D
  // addr->R14
@0
D=A
@LCL
D=D+M
@R14
M=D
  // R13->addrOfR14
@R13
D=M
@R14
A=M
M=D
// [END] pop local 0
// push local 1(type 1)
@1
D=A
@LCL
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push local 1
// push constant 1(type 1)
@1
D=A
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push constant 1
// add(type 0)
@BACK.5
D=A
@BACK
M=D
@ADD
0;JMP
(BACK.5)
// [END] add
// pop local 1(type 2)
  //POP
@SP
M=M-1
A=M
D=M
  // D->R13
@R13
M=D
  // addr->R14
@1
D=A
@LCL
D=D+M
@R14
M=D
  // R13->addrOfR14
@R13
D=M
@R14
A=M
M=D
// [END] pop local 1
// goto WHILE_LOOP(type 4)
  // goto
@Test3.vm.mult$WHILE_LOOP
0;JMP
// [END] goto WHILE_LOOP
// label WHILE_END(type 3)
  // label
(Test3.vm.mult$WHILE_END)
// [END] label WHILE_END
// push local 0(type 1)
@0
D=A
@LCL
A=D+M
D=M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
// [END] push local 0
// return(type 7)
  // return
@LCL
D=M
@R13
M=D
@5
D=A
@R13
D=M-D
A=D
D=M
@R14
M=D
  //POP
@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D
@ARG
D=M+1
@SP
M=D
@R13
D=M-1
A=D
D=M
@THAT
M=D
@2
D=A
@R13
D=M-D
A=D
D=M
@THIS
M=D
@3
D=A
@R13
D=M-D
A=D
D=M
@ARG
M=D
@4
D=A
@R13
D=M-D
A=D
D=M
@LCL
M=D
@R14
A=M
0;JMP
// [END] return
@END
0;JMP
//
// PROCEDURES SECTION
//
(ADD)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M
@R13
D=D+M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(SUB)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M
@R13
D=D-M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(NEG)
  //POP
@SP
M=M-1
A=M
D=M
D=-D
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(GT)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M

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
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(LT)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M

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
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(EQ)
  //POP
@SP
M=M-1
A=M
D=M
  // save to R13
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M
  // compute
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
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(AND)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M
@R13
D=D&M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(OR)
  //POP
@SP
M=M-1
A=M
D=M
@R13
M=D
  //POP
@SP
M=M-1
A=M
D=M
@R13
D=D|M
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(NOT)
  //POP
@SP
M=M-1
A=M
D=M
D=!D
  //PUSH
@SP
A=M
M=D
@SP
M=M+1
@DONE
0;JMP
(DONE)
@BACK
A=M
0;JMP
(END)
@END
0;JMP
