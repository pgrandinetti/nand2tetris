// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/4/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed,
// the screen should be cleared.


(LOOP)
    @KBD
    D=M
    @SET0
    D;JEQ
    @SET1
    0;JMP
(DONE)
    // update 'prev', then continue LOOP
    @KBD
    D=M
    @prev
    M=D
    @LOOP
    0;JMP

(SET1)
    @set
    M=-1
    @SET
    0;JMP

(SET0)
    @set
    M=0
    @SET
    0;JMP

(SET)
    // fills SCREEN x8K with the value in register 'set'
    @16384
    D=A
    @start
    M=D
    @24576
    D=A
    @end
    M=D

(SCREENLOOP)
    // if start == end goto main LOOP
    @start
    D=M
    @end
    D=M-D
    @DONE
    D;JEQ

    // fill current RAM[start] with 'set'
    @set
    D=M
    @start
    A=M
    M=D

    // increment start
    @start
    M=M+1
    @SCREENLOOP
    0;JMP
