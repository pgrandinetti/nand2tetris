package main

import "testing"

func TestAdvance(t *testing.T) {
	inpt := "   push    local 2  "
	res := Advance(inpt)
	exp := "push local 2"
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

	inpt = "  // comment  "
	res = Advance(inpt)
	exp = ""
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

    inpt = "not  // comment"
    res = Advance(inpt)
    exp = "not"
    if res != exp {
        t.Errorf("expected |%s| computed |%s|", exp, res)
    }

    inpt = "push local 1 //"
    res = Advance(inpt)
    exp = "push local 1"
    if res != exp {
        t.Errorf("expected |%s| computed |%s|", exp, res)
    }

    inpt = "push argument 1"
    res = Advance(inpt)
    exp = "push argument 1"
    if res != exp {
        t.Errorf("expected |%s| computed |%s|", exp, res)
    }
}

func TestType(t *testing.T) {
    instr := "call mult 2"
    exp := C_CALL
    res := CommandType(instr)
	if res != exp {
		t.Errorf("expected |%d| computed |%d|", exp, res)
	}

    instr = "push local 2"
    exp = C_PUSH
    res = CommandType(instr)
    if res != exp {
        t.Errorf("expected |%d| computed |%d|", exp, res)
    }

    instr = "label END"
    exp = C_LABEL
    res = CommandType(instr)
    if res != exp {
        t.Errorf("expected |%d| computed |%d|", exp, res)
    }
}

func TestArg1(t *testing.T) {
	instr := "push local 2"
	exp := "local"
	res := Arg1(instr, C_PUSH)
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

	instr = "add"
	exp = instr
	res = Arg1(instr, C_ARITHMETIC)
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

	instr = "pop x"
	exp = "x"
	res = Arg1(instr, C_POP)
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

	instr = "function foo 5"
	exp = "foo"
	res = Arg1(instr, C_FUNCTION)
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}

    instr = "push constant 2"
    exp = "constant"
	res = Arg1(instr, C_PUSH)
	if res != exp {
		t.Errorf("expected |%s| computed |%s|", exp, res)
	}
}

func TestArg2(t *testing.T) {
	instr := "push local 2"
	exp := 2
	res := Arg2(instr)
	if res != exp {
		t.Errorf("expected |%d| computed |%d|", exp, res)
	}

	instr = "function foo 5"
	exp = 5
	res = Arg2(instr)
	if res != exp {
		t.Errorf("expected |%d| computed |%d|", exp, res)
	}

}
