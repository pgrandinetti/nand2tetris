package main

import (
	"fmt"
)

var labelCount int = 0

func process(tkn Token, s string, tknType TokenT) string {
	// Process the terminal nodes of the tree.
	// Idea is to produce VM code that adds one (and only one) element to the stack.
	if tkn != s {
		panic(fmt.Sprintf("expected |%s| received token |%s|", s, tkn))
	}
	if TokenType(tkn) != tknType {
		panic(fmt.Sprintf("expected %v token is %s with type %v", tknType, tkn, TokenType(tkn)))
	}
	if tknType == SYMBOL {
		if tkn == "+" {
			return "add\n"
		}
		if tkn == "-" {
			return "sub\n"
		}
		if tkn == "*" {
			return "call Math.multiply 2\n"
		}
		if tkn == "/" {
			return "call Math.divide 2\n"
		}
		if tkn == "&" {
			return "and\n"
		}
		if tkn == "|" {
			return "or\n"
		}
		if tkn == ">" {
			return "gt\n"
		}
		if tkn == "<" {
			return "lt\n"
		}
		if tkn == "=" {
			return "eq\n"
		}
		// UnaryOps are handled in another func
		// else is not an operator
		return ""
	} else if tknType == KEYWORD {
		if tkn == "true" {
			return "push constant 0\nnot\n"
		}
		if tkn == "false" || tkn == "null" {
			return "push constant 0\n"
		}
		if tkn == "this" {
			return "push pointer 0\n"
		}
		return ""
	} else if tknType == INT_CONST {
		return fmt.Sprintf("push constant %s\n", tkn)
	} else if tknType == STRING_CONST {
		// String API operates on 'this' object
		// hence save the current value before process begins
		st := "push pointer 0\npop temp 0\n"                                   // 'this' will be set to 'temp 0' before exiting
		st += fmt.Sprintf("push constant %d\ncall String.new 1\n", len(tkn)-2) // minus twice "
		for i := 1; i < len(tkn)-1; i++ {
			// String.new (and appendChar too) puts a new object pointer on the top of the stack
			st += fmt.Sprintf("push constant %d\n", tkn[i])
			st += "call String.appendChar 2\n"
		}
		// restore 'this 0' and leave the new String at the top of the stack
		st += "push temp 0\npop pointer 0\n"
		return st
	} else {
		// IDENTIFIER
		kind, err := KindOf(tkn)
		if err != nil {
			// must be a function's or method's name
			return ""
			//panic(fmt.Sprintf("identifier %s not found in symbol table", tkn))
		}
		idx, _ := IndexOf(tkn)
		if kind == STATICV {
			return fmt.Sprintf("push static %d\n", idx)
		}
		if kind == FIELDV {
			return fmt.Sprintf("push this %d\n", idx)
		}
		if kind == ARGV {
			return fmt.Sprintf("push argument %d\n", idx)
		}
		if kind == VARV {
			return fmt.Sprintf("push local %d\n", idx)
		}
		panic(fmt.Sprintf("kind %+v not valid for identifier %s", kind, tkn))
	}
}

func CompileExpression(content []byte, next int) (int, string) {
	var prev int
	var err error
	var tkn Token

	var st string       // final VM code
	var st1, st2 string // temp vars

	next, st = CompileTerm(content, next)
	prev = next
	tkn, next, err = Advance(content, next)
	for err == nil {
		if tkn == "+" || tkn == "-" || tkn == "*" || tkn == "/" || tkn == "&" || tkn == "|" || tkn == ">" || tkn == "<" || tkn == "=" {
			next, st1 = CompileOp(content, prev)
			next, st2 = CompileTerm(content, next)
			prev = next
			tkn, next, err = Advance(content, next)
			st += st2 + st1 // postfix
		} else {
			break
		}
	}
	return prev, st
}

func CompileTerm(content []byte, next int) (int, string) {
	var tkn, tkn2 Token
	var current int
	var st string       // VM code to return
	var st1, st2 string // temp vars

	current = next
	tkn, next, _ = Advance(content, next)
	if tkn == "true" || tkn == "false" || tkn == "null" || tkn == "this" {
		st = process(tkn, tkn, KEYWORD)
	} else if tkn[0] >= '0' && tkn[0] <= '9' {
		st = process(tkn, tkn, INT_CONST)
	} else if tkn[0] == '"' {
		st = process(tkn, tkn, STRING_CONST)
	} else if tkn == "(" {
		process(tkn, tkn, SYMBOL)
		next, st = CompileExpression(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
	} else if tkn == "-" || tkn == "~" {
		next, st1 = CompileUnaryOp(content, current)
		next, st2 = CompileTerm(content, next)
		st = st2 + st1 // postfix
	} else {
		// need to look ahead to tell 'varName' | 'varName[expression]' | 'subroutineCall'
		tkn2, next, _ = Advance(content, next)
		if tkn2 == "[" {
			next, st1 = CompileVarName(content, current)
			tkn, next, _ = Advance(content, next)
			process(tkn, "[", SYMBOL)
			next, st2 = CompileExpression(content, next)
			tkn, next, _ = Advance(content, next)
			process(tkn, "]", SYMBOL)
			st = st1 + st2 + "add\n" // (arr+expr)
			st += "pop pointer 1\n"  // align THAT
			st += "push that 0\n"    // puts on the stack *(arr+expr)
		} else if tkn2 == "." || tkn2 == "(" {
			next, st = CompileSubroutineCall(content, current)
		} else {
			next, st = CompileVarName(content, current)
		}
	}
	return next, st
}

func CompileUnaryOp(content []byte, next int) (int, string) {
	var st string
	tkn, next, _ := Advance(content, next)
	if tkn == "-" {
		st = "neg\n"
	} else if tkn == "~" {
		st = "not\n"
	}
	return next, st
}

func CompileOp(content []byte, next int) (int, string) {
	tkn, next, _ := Advance(content, next)
	st := process(tkn, tkn, SYMBOL)
	return next, st
}

func compileIdentifier(content []byte, next int) (int, string) {
	tkn, next, _ := Advance(content, next)
	st := process(tkn, tkn, IDENTIFIER)
	return next, st
}

func CompileClassName(content []byte, next int) (int, string) {
	return compileIdentifier(content, next)
}

func CompileSubroutineName(content []byte, next int) (int, string) {
	return compileIdentifier(content, next)
}

func CompileVarName(content []byte, next int) (int, string) {
	return compileIdentifier(content, next)
}

func CompileSubroutineCall(content []byte, next int) (int, string) {
	var tkn Token
	var st, st2, varName, funcName string // VM code, temp vars
	var nExprs int

	tkn, next, _ = Advance(content, next)
	st = process(tkn, tkn, IDENTIFIER) // can be both varName and funcName
	varName = tkn
	funcName = tkn
	tkn, next, _ = Advance(content, next)
	if tkn == "." {
		process(tkn, tkn, SYMBOL)
		tkn, next, _ = Advance(content, next)
		process(tkn, tkn, IDENTIFIER)
		funcName = tkn
		tkn, next, _ = Advance(content, next)
		process(tkn, "(", SYMBOL)
		next, st2, nExprs = CompileExpressionList(content, next)
		st += st2
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
		jackType, _ := TypeOf(varName)
		if jackType == "" {
			// then it's a class name
			st += fmt.Sprintf("call %s.%s %d\n", varName, funcName, nExprs)
		} else {
			st += fmt.Sprintf("call %s.%s %d\n", jackType, funcName, nExprs+1)
		}
	} else if tkn == "(" {
		// this must be a method called onto 'this' object
		varName = "this"
		st = "push pointer 0\n" + st
		process(tkn, tkn, SYMBOL)
		next, st2, nExprs = CompileExpressionList(content, next)
		st += st2
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
		//jackType, _ := TypeOf(varName)
		//st += fmt.Sprintf("call %s.%s %d\n", jackType, funcName, nExprs+1)
		st += fmt.Sprintf("call %s.%s %d\n", CurrentClass, funcName, nExprs+1)
	} else {
		panic(fmt.Sprintf("unterminated subroutineCall at index %d. expected |.| or |(|", next))
	}
	return next, st
}

func CompileExpressionList(content []byte, next int) (int, string, int) {
	var tkn Token
	var current, nExprs int
	var st, st2 string // VM code, temp var

	current = next
	tkn, next, _ = Advance(content, next)
	if tkn == ")" {
		// expressionList is empty
		return current, "", 0
	}
	nExprs++
	next, st2 = CompileExpression(content, current)
	st += st2
	current = next
	tkn, next, _ = Advance(content, next)
	for tkn == "," {
		nExprs++
		process(tkn, tkn, SYMBOL)
		next, st2 = CompileExpression(content, next)
		st += st2
		current = next
		tkn, next, _ = Advance(content, next)
	}
	return current, st, nExprs
}

func CompileClassVarDec(content []byte, next int) (int, string) {
	var tkn Token
	var err error

	// Used for the SymbolTable.
	var kind VarKind
	var jackType, name string

	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, KEYWORD) // static | field
	if tkn == "static" {
		kind = STATICV
	} else {
		kind = FIELDV
	}

	jackType, _, _ = Advance(content, next)
	next, _ = CompileType(content, next)

	name, _, _ = Advance(content, next)
	next, _ = CompileVarName(content, next)

	Define(name, jackType, kind)

	for {
		tkn, next, err = Advance(content, next)
		if err != nil {
			panic("incomplete 'classVarDec'. expected ';'")
		}
		if tkn == "," {
			process(tkn, ",", SYMBOL)
			name, _, _ = Advance(content, next)
			next, _ = CompileVarName(content, next)
			Define(name, jackType, kind)
		} else {
			process(tkn, ";", SYMBOL)
			break
		}
	}
	return next, ""
}

func CompileType(content []byte, next int) (int, string) {
	var tkn Token
	tkn, next, _ = Advance(content, next)
	if tkn == "int" || tkn == "char" || tkn == "boolean" {
		process(tkn, tkn, KEYWORD)
	} else {
		process(tkn, tkn, IDENTIFIER)
	}
	return next, ""
}

func CompileSubroutineDec(content []byte, next int) (int, string) {
	var tkn Token
	var st, st2 string // VM code

	var isMethod bool = false
	var isConstr bool = false
	var funcName string

	ResetRoutine() // Symbol Table

	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, KEYWORD) // constructor | method | function
	if tkn == "method" {
		isMethod = true
		DefineThis()
	}
	if tkn == "constructor" {
		isConstr = true
	}

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn == "void" {
		process(tkn, tkn, KEYWORD)
	} else {
		next, _ = CompileType(content, prev)
	}

	funcName, _, _ = Advance(content, next) // get the routine name. do not advance pointer.
	next, _ = CompileSubroutineName(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next, _ = CompileParameterList(content, next) // must populate the routine's SymbolTable

	tkn, next, _ = Advance(content, next)
	process(tkn, ")", SYMBOL)
	next, st2 = CompileSubroutineBody(content, next)

	st = fmt.Sprintf("function %s.%s %d\n", CurrentClass, funcName, VarCount(VARV))
	if isMethod {
		st += "push argument 0\npop pointer 0\n"
	} else if isConstr {
		nFields := VarCount(FIELDV)
		st += fmt.Sprintf("push constant %d\n", nFields)
		st += "call Memory.alloc 1\n"
		st += "pop pointer 0\n"
	}

	st += st2
	// Book says to add the two lines below
	// but the syntax enforces 'return this;' at the end of any constructor
	// and 'return this;' is already compiled as 'push this 0; return;' by CompileReturn.
	//if isConstr {
	//    st += "push pointer 0\nreturn\n"
	//}

	return next, st
}

func CompileParameterList(content []byte, next int) (int, string) {
	var tkn Token
	var err error

	// To use with SymbolTable
	var jackType, name string

	prev := next
	tkn, _, err = Advance(content, next)
	if tkn == ")" {
		return prev, ""
	}
	jackType, _, _ = Advance(content, prev)
	next, _ = CompileType(content, prev)
	name, _, _ = Advance(content, next)
	next, _ = CompileVarName(content, next)
	Define(name, jackType, ARGV)

	prev = next
	tkn, next, err = Advance(content, next)
	for tkn == "," && err == nil {
		process(tkn, tkn, SYMBOL) // ','
		jackType, _, _ = Advance(content, next)
		next, _ = CompileType(content, next)
		name, _, _ = Advance(content, next)
		next, _ = CompileVarName(content, next)
		Define(name, jackType, ARGV)
		prev = next
		tkn, next, err = Advance(content, next)
	}
	if err != nil {
		panic("unterminated parameterList: expected |)|")
	}
	return prev, ""
}

func CompileReturn(content []byte, next int) (int, string) {
	var tkn Token
	var st string // VM code

	tkn, next, _ = Advance(content, next)
	process(tkn, "return", KEYWORD)

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn != ";" {
		next, st = CompileExpression(content, prev)
		tkn, next, _ = Advance(content, next)
		process(tkn, ";", SYMBOL)
	} else {
		st = ""
		process(tkn, tkn, SYMBOL)
	}
	return next, st + "return\n"
}

func CompileLet(content []byte, next int) (int, string) {
	var tkn Token
	var st, st2, stExpr string // VM code
	var varName string

	tkn, next, _ = Advance(content, next)
	process(tkn, "let", KEYWORD)
	varName, _, _ = Advance(content, next)
	next, _ = CompileVarName(content, next)

	tkn, next, _ = Advance(content, next)
	if tkn == "[" {
		process(tkn, tkn, SYMBOL)
		next, stExpr = CompileExpression(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, "]", SYMBOL)
		tkn, next, _ = Advance(content, next)
	}
	process(tkn, "=", SYMBOL)
	next, st2 = CompileExpression(content, next)
	st = st2
	tkn, next, _ = Advance(content, next)
	process(tkn, ";", SYMBOL)

	kind, _ := KindOf(varName)
	idx, _ := IndexOf(varName)
	if stExpr == "" {
		// simple VarName (no array)
		if kind == ARGV {
			st += fmt.Sprintf("pop argument %d\n", idx)
		} else if kind == STATICV {
			st += fmt.Sprintf("pop static %d\n", idx)
		} else if kind == FIELDV {
			st += fmt.Sprintf("pop this %d\n", idx)
		} else { // VARV
			st += fmt.Sprintf("pop local %d\n", idx)
		}
	} else {
		// handle array 'let arr[expr] = expr2'
		st += stExpr
		if kind == ARGV {
			st += fmt.Sprintf("push argument %d\n", idx)
		} else if kind == STATICV {
			st += fmt.Sprintf("push static %d\n", idx)
		} else if kind == FIELDV {
			st += fmt.Sprintf("push this %d\n", idx)
		} else { // VARV
			st += fmt.Sprintf("push local %d\n", idx)
		}
		st += "add\n"           // compute (arr + expr)
		st += "pop pointer 1\n" // align THAT
		// now the stack only has the result of expr2
		st += "pop that 0\n"
	}
	return next, st
}

func CompileDo(content []byte, next int) (int, string) {
	var tkn Token
	var st string // VM code

	tkn, next, _ = Advance(content, next)
	process(tkn, "do", KEYWORD)
	next, st = CompileSubroutineCall(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ";", SYMBOL)
	st += "pop temp 0\n"
	return next, st
}

func CompileIf(content []byte, next int) (int, string) {
	var tkn Token
	var stExpr, stIf, stElse, st string // VM code

	labelL1 := fmt.Sprintf("%s.IF.%d", CurrentClass, labelCount)
	labelCount++
	labelL2 := fmt.Sprintf("%s.ELSE.%d", CurrentClass, labelCount)
	labelCount++

	tkn, next, _ = Advance(content, next)
	process(tkn, "if", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next, stExpr = CompileExpression(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ")", SYMBOL)
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)
	next, stIf = CompileStatements(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn == "else" {
		process(tkn, tkn, KEYWORD)
		tkn, next, _ = Advance(content, next)
		process(tkn, "{", SYMBOL)
		next, stElse = CompileStatements(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, "}", SYMBOL)
		prev = next
	}

	st = stExpr + fmt.Sprintf("not\nif-goto %s\n", labelL1)
	st += stIf
	st += fmt.Sprintf("goto %s\n", labelL2)
	st += fmt.Sprintf("label %s\n", labelL1)
	st += stElse
	st += fmt.Sprintf("label %s\n", labelL2)

	return prev, st
}

func CompileWhile(content []byte, next int) (int, string) {
	var tkn Token
	var stExpr, stStmts, st string // VM code

	label1 := fmt.Sprintf("%s.WHILE.%d", CurrentClass, labelCount)
	labelCount++
	label2 := fmt.Sprintf("%s.WHILE.END.%d", CurrentClass, labelCount)
	labelCount++

	tkn, next, _ = Advance(content, next)
	process(tkn, "while", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next, stExpr = CompileExpression(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ")", SYMBOL)
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)
	next, stStmts = CompileStatements(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)

	st = "label " + label1 + "\n" + stExpr + "not\n"
	st += fmt.Sprintf("if-goto %s\n", label2)
	st += stStmts
	st += fmt.Sprintf("goto %s\n", label1)
	st += "label " + label2 + "\n"

	return next, st
}

func CompileStatements(content []byte, next int) (int, string) {
	var tkn Token
	var prev int
	var err error
	var st, st2 string // VM code
	st = ""

	prev = next
	tkn, next, err = Advance(content, next)
	for err == nil {
		if tkn == "let" {
			next, st2 = CompileLet(content, prev)
		} else if tkn == "if" {
			next, st2 = CompileIf(content, prev)
		} else if tkn == "while" {
			next, st2 = CompileWhile(content, prev)
		} else if tkn == "do" {
			next, st2 = CompileDo(content, prev)
		} else if tkn == "return" {
			next, st2 = CompileReturn(content, prev)
		} else {
			return prev, st
		}
		st += st2
		prev = next
		tkn, next, err = Advance(content, next)
	}
	panic("unterminated statements")
}

func CompileVarDec(content []byte, next int) (int, string) {
	var tkn Token

	// Used with SymbolTable.
	var jackType, name string

	tkn, next, _ = Advance(content, next)
	process(tkn, "var", KEYWORD)
	jackType, _, _ = Advance(content, next)
	next, _ = CompileType(content, next)
	name, _, _ = Advance(content, next)
	next, _ = CompileVarName(content, next)
	Define(name, jackType, VARV)

	tkn, next, _ = Advance(content, next)
	for tkn == "," {
		process(tkn, tkn, SYMBOL)
		name, _, _ = Advance(content, next)
		next, _ = CompileVarName(content, next)
		Define(name, jackType, VARV)
		tkn, next, _ = Advance(content, next)
	}
	process(tkn, ";", SYMBOL)
	return next, ""
}

func CompileSubroutineBody(content []byte, next int) (int, string) {
	var tkn Token
	var prev int
	var err error
	var st string // VM code

	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)

	prev = next
	tkn, next, _ = Advance(content, next)
	for tkn == "var" {
		next, _ = CompileVarDec(content, prev)
		prev = next
		tkn, next, err = Advance(content, next)
		if err != nil {
			panic("unterminated subroutineBody")
		}
	}
	next, st = CompileStatements(content, prev)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)
	return next, st
}

func CompileClass(content []byte, next int) (int, string) {
	var tkn Token
	var current int
	var st, st2 string // VM code

	ResetClass() // for SymbolTable

	tkn, next, _ = Advance(content, next)
	process(tkn, "class", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, IDENTIFIER)

	CurrentClass = tkn // SymbolTable

	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)

	current = next
	tkn, next, _ = Advance(content, next)
	for tkn == "static" || tkn == "field" {
		next, _ = CompileClassVarDec(content, current)
		current = next
		tkn, next, _ = Advance(content, next)
	}
	for tkn == "constructor" || tkn == "function" || tkn == "method" {
		next, st2 = CompileSubroutineDec(content, current)
		st += st2
		current = next
		tkn, next, _ = Advance(content, next)
		//fmt.Printf("ROUTINE TABLE: %+v\n", RoutineVars)
	}
	process(tkn, "}", SYMBOL)
	return next, st
}

/*
// Test that can be executed manually.
func main() {
    content := `
        class Point {
            field int x, y;
            static int totCount;
            constructor Point new (int ax, int ay) {
                let x = ax;
                let y = ay;
                let totCount = totCount + 1;
                return this;
            }
            method int sum (Point other) {
                var int newx, newy;
                let newx = x + other.getx();
                let newy = y + other.gety();
                return Point.new(newx, newy);
            }
            method int getx () { return x; }
            method int gety() { return y; }
            function int getCount() {
                return totCount;
            }
        }
    `
    fmt.Println(content)
    _, code := CompileClass([]byte(content), 0)
    fmt.Println(code)

    content = `
        class Main {
            function void main () {
                var Array x;
                var Point p1, p2, p3;
                var int i;

                let x = Array.new(3);
                let i = 0;
                let p1 = Point.new(7, 3);

                while (i < 10) {
                    if (i < 5) {
                        let x[1] = x[1] + i;  // 1+2+3+4 = 10
                    } else {
                        let x[1] = x[1] + 1; // (10+)1+1+1+1+1 = 15
                    }
                    let i = i + 1;
                }

                let p2 = Point.new(3, x[1]);  // (3, 15)
                let p3 = p1.sum(p2);  // (10, 18)

                do Output.printString("P3.x is ");
                do Output.printInt(p3.getx());
                do Output.println();
                do Output.printString("y is ");
                do Output.printInt(p3.gety());
                do Output.println();
            }
        }
    `
    fmt.Println(content)
    _, code = CompileClass([]byte(content), 0)
    fmt.Println(code)
}
*/
