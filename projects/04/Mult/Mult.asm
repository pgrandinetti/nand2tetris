// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
// The algorithm is based on repetitive addition.


    // i = 0; operation counter
    @i
    M=0

    // result = 0;
    @result
    M=0

(LOOP)
    // if i == R1 goto DONE
    @i
    D=M
    @R1
    D=D-M
    @DONE
    D;JEQ

    // result += R0
    @result
    D=M
    @R0
    D=D+M
    @result
    M=D

    // i += 1
    @i
    M=M+1

    // repeat
    @LOOP
    0;JMP

(DONE)
    @result
    D=M
    @R2
    M=D
    @END
    0;JMP

(END)
    @END
    0;JMP
