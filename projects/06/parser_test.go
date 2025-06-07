package main

import "testing"

func TestAdvance(t *testing.T) {
	inpt := "  x t  "
	res := Advance(inpt)
	exp := "x t"
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	inpt = "// abfng"
	res = Advance(inpt)
	exp = ""
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	inpt = "  // abfng"
	res = Advance(inpt)
	exp = ""
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	inpt = "  "
	res = Advance(inpt)
	exp = ""
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}
}

func TestInstrType(t *testing.T) {
	instr := "@7"
	res := InstructionType(instr)
	exp := A_INSTRUCTION
	if res != exp {
		t.Errorf("expected |%d|. computed |%d|", exp, res)
	}

	instr = "DM=D+1;JTE"
	res = InstructionType(instr)
	exp = C_INSTRUCTION
	if res != exp {
		t.Errorf("expected |%d|. computed |%d|", exp, res)
	}

	instr = "(LABEL)"
	res = InstructionType(instr)
	exp = L_INSTRUCTION
	if res != exp {
		t.Errorf("expected |%d|. computed |%d|", exp, res)
	}
}

func TestSymbol(t *testing.T) {
	instr := "@123"
	exp := "123"
	res := Symbol(instr, A_INSTRUCTION)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "@var1"
	exp = "var1"
	res = Symbol(instr, A_INSTRUCTION)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "(LABEL)"
	exp = "LABEL"
	res = Symbol(instr, L_INSTRUCTION)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}
}

func TestDest(t *testing.T) {
	instr := "M=D+1;JEQ"
	exp := "M"
	res := Dest(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "DM=D+1;JMP"
	exp = "DM"
	res = Dest(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "D&A;JMP"
	exp = ""
	res = Dest(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}
}

func TestComp(t *testing.T) {
	instr := "M=D+1;JEQ"
	exp := "D+1"
	res := Comp(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "DM=A|M;JMP"
	exp = "A|M"
	res = Comp(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "D&A;JMP"
	exp = "D&A"
	res = Comp(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "D&A"
	exp = "D&A"
	res = Comp(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "D=D&A"
	exp = "D&A"
	res = Comp(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}
}

func TestJump(t *testing.T) {
	instr := "M=D+1;JEQ"
	exp := "JEQ"
	res := Jump(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "DM=A|M;JMP"
	exp = "JMP"
	res = Jump(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}

	instr = "DM=A|M"
	exp = ""
	res = Jump(instr)
	if res != exp {
		t.Errorf("expected |%s|. computed |%s|", exp, res)
	}
}
