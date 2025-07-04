package main

import "fmt"

func process(tkn Token, s string, tknType TokenT) {
	if tkn != s {
		panic(fmt.Sprintf("expected |%s| received token |%s|", s, tkn))
	}
	if TokenType(tkn) != tknType {
		panic(fmt.Sprintf("expected %v token is %s with type %v", tknType, tkn, TokenType(tkn)))
	}
	if tknType == SYMBOL {
		s1 := s
		if s == "<" {
			s1 = "&lt;"
		} else if s == ">" {
			s1 = "&gt;"
		} else if s == "&" {
			s1 = "&amp;"
		}
		fmt.Printf("<symbol> %s </symbol>\n", s1)
	} else if tknType == KEYWORD {
		fmt.Printf("<keyword> %s </keyword>\n", s)
	} else if tknType == INT_CONST {
		fmt.Printf("<integerConstant> %s </integerConstant>\n", s)
	} else if tknType == STRING_CONST {
		fmt.Printf("<stringConstant> %s </stringConstant>\n", s[1:len(s)-1]) // without "s
	} else {
		fmt.Printf("<identifier> %s </identifier>\n", s)
	}
}

func CompileClass(content []byte, next int) int {
	var tkn Token
	var current int
	fmt.Println("<class>")
	tkn, next, _ = Advance(content, next)
	process(tkn, "class", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, IDENTIFIER)
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)

	current = next
	tkn, next, _ = Advance(content, next)
	for tkn == "static" || tkn == "field" {
		next = CompileClassVarDec(content, current)
		current = next
		tkn, next, _ = Advance(content, next)
	}
	for tkn == "constructor" || tkn == "function" || tkn == "method" {
		next = CompileSubroutineDec(content, current)
		current = next
		tkn, next, _ = Advance(content, next)
	}
	process(tkn, "}", SYMBOL)
	fmt.Println("</class>")
	return next
}

func CompileClassVarDec(content []byte, next int) int {
	var tkn Token
	var err error

	fmt.Println("<classVarDec>")
	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, KEYWORD) // static | field
	next = CompileType(content, next)
	next = CompileVarName(content, next)
	for {
		tkn, next, err = Advance(content, next)
		if err != nil {
			panic("incomplete 'classVarDec'. expected ';'")
		}
		if tkn == "," {
			process(tkn, ",", SYMBOL)
			next = CompileVarName(content, next)
		} else {
			process(tkn, ";", SYMBOL)
			break
		}
	}
	fmt.Println("</classVarDec>")
	return next
}

func CompileType(content []byte, next int) int {
	var tkn Token
	tkn, next, _ = Advance(content, next)
	if tkn == "int" || tkn == "char" || tkn == "boolean" {
		process(tkn, tkn, KEYWORD)
	} else {
		process(tkn, tkn, IDENTIFIER)
	}
	return next
}

func CompileSubroutineDec(content []byte, next int) int {
	var tkn Token
	fmt.Println("<subroutineDec>")

	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, KEYWORD) // constructor | method | function

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn == "void" {
		process(tkn, tkn, KEYWORD)
	} else {
		next = CompileType(content, prev)
	}

	next = CompileSubroutineName(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next = CompileParameterList(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ")", SYMBOL)
	next = CompileSubroutineBody(content, next)

	fmt.Println("</subroutineDec>")
	return next
}

func CompileParameterList(content []byte, next int) int {
	fmt.Println("<parameterList>")
	var tkn Token
	var err error
	prev := next
	tkn, _, err = Advance(content, next)
	if tkn == ")" {
		// no parameter list
		fmt.Println("</parameterList>")
		return prev
	}
	next = CompileType(content, prev)
	next = CompileVarName(content, next)

	prev = next
	tkn, next, err = Advance(content, next)
	for tkn == "," && err == nil {
		process(tkn, tkn, SYMBOL) // ','
		next = CompileType(content, next)
		next = CompileVarName(content, next)
		prev = next
		tkn, next, err = Advance(content, next)
	}
	if err != nil {
		panic("unterminated parameterList: expected |)|")
	}
	return prev
}

func CompileSubroutineBody(content []byte, next int) int {
	fmt.Println("<subroutineBody>")
	var tkn Token
	var prev int
	var err error
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)

	prev = next
	tkn, next, _ = Advance(content, next)
	for tkn == "var" {
		next = CompileVarDec(content, prev)
		prev = next
		tkn, next, err = Advance(content, next)
		if err != nil {
			panic("unterminated subroutineBody")
		}
	}
	next = CompileStatements(content, prev)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)
	fmt.Println("</subroutineBody>")
	return next
}

func CompileVarDec(content []byte, next int) int {
	fmt.Println("<varDec>")
	var tkn Token
	tkn, next, _ = Advance(content, next)
	process(tkn, "var", KEYWORD)
	next = CompileType(content, next)
	next = CompileVarName(content, next)

	tkn, next, _ = Advance(content, next)
	for tkn == "," {
		process(tkn, tkn, SYMBOL)
		next = CompileVarName(content, next)
		tkn, next, _ = Advance(content, next)
	}
	process(tkn, ";", SYMBOL)
	fmt.Println("</varDec>")
	return next
}

func compileIdentifier(content []byte, next int) int {
	tkn, next, _ := Advance(content, next)
	process(tkn, tkn, IDENTIFIER)
	return next
}

func CompileClassName(content []byte, next int) int {
	return compileIdentifier(content, next)
}

func CompileSubroutineName(content []byte, next int) int {
	return compileIdentifier(content, next)
}

func CompileVarName(content []byte, next int) int {
	return compileIdentifier(content, next)
}

func CompileStatements(content []byte, next int) int {
	var tkn Token
	var prev int
	var err error

	prev = next
	tkn, next, err = Advance(content, next)
	for err == nil {
		if tkn == "let" {
			next = CompileLet(content, prev)
		} else if tkn == "if" {
			next = CompileIf(content, prev)
		} else if tkn == "while" {
			next = CompileWhile(content, prev)
		} else if tkn == "do" {
			next = CompileDo(content, prev)
		} else if tkn == "return" {
			next = CompileReturn(content, prev)
		} else {
			// no statement
			return prev
		}
		prev = next
		tkn, next, err = Advance(content, next)
	}
	panic("unterminated statements")
}

func CompileLet(content []byte, next int) int {
	fmt.Println("<letStatement>")
	var tkn Token

	tkn, next, _ = Advance(content, next)
	process(tkn, "let", KEYWORD)
	next = CompileVarName(content, next)

	tkn, next, _ = Advance(content, next)
	if tkn == "[" {
		process(tkn, tkn, SYMBOL)
		next = CompileExpression(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, "]", SYMBOL)
		tkn, next, _ = Advance(content, next)
	}
	process(tkn, "=", SYMBOL)
	next = CompileExpression(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ";", SYMBOL)
	fmt.Println("</letStatement>")
	return next
}

func CompileIf(content []byte, next int) int {
	fmt.Println("<ifStatement>")
	var tkn Token

	tkn, next, _ = Advance(content, next)
	process(tkn, "if", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next = CompileExpression(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ")", SYMBOL)
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)
	next = CompileStatements(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn == "else" {
		process(tkn, tkn, KEYWORD)
		tkn, next, _ = Advance(content, next)
		process(tkn, "{", SYMBOL)
		next = CompileStatements(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, "}", SYMBOL)
		prev = next
	}
	fmt.Println("</ifStatement>")
	return prev
}

func CompileWhile(content []byte, next int) int {
	var tkn Token
	fmt.Println("<whileStatement>")
	tkn, next, _ = Advance(content, next)
	process(tkn, "while", KEYWORD)
	tkn, next, _ = Advance(content, next)
	process(tkn, "(", SYMBOL)
	next = CompileExpression(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "{", SYMBOL)
	next = CompileStatements(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, "}", SYMBOL)
	fmt.Println("</whileStatement")
	return next
}

func CompileDo(content []byte, next int) int {
	fmt.Println("<doStatement>")
	var tkn Token
	tkn, next, _ = Advance(content, next)
	process(tkn, "do", KEYWORD)
	next = CompileSubroutineCall(content, next)
	tkn, next, _ = Advance(content, next)
	process(tkn, ";", SYMBOL)
	fmt.Println("</doStatement>")
	return next
}

func CompileReturn(content []byte, next int) int {
	fmt.Println("<returnStatement>")
	var tkn Token
	tkn, next, _ = Advance(content, next)
	process(tkn, "return", KEYWORD)

	prev := next
	tkn, next, _ = Advance(content, next)
	if tkn != ";" {
		next = CompileExpression(content, prev)
		tkn, next, _ = Advance(content, next)
		process(tkn, ";", SYMBOL)
	} else {
		process(tkn, tkn, SYMBOL)
	}
	fmt.Println("</returnStatement>")
	return next
}

func CompileExpression(content []byte, next int) int {
	fmt.Println("<expression>")
	var prev int
	var err error
	var tkn Token

	next = CompileTerm(content, next)
	prev = next
	tkn, next, err = Advance(content, next)
	for err == nil {
		if tkn == "+" || tkn == "-" || tkn == "*" || tkn == "/" || tkn == "&" || tkn == "|" || tkn == ">" || tkn == "<" || tkn == "=" {
			next = CompileOp(content, prev)
			next = CompileTerm(content, next)
			prev = next
			tkn, next, err = Advance(content, next)
		} else {
			break
		}
	}
	fmt.Println("</expression>")
	return prev
}

func CompileTerm(content []byte, next int) int {
	fmt.Println("<term>")
	var tkn, tkn2 Token
	var current int

	current = next
	tkn, next, _ = Advance(content, next)
	if tkn == "true" || tkn == "false" || tkn == "null" || tkn == "this" {
		process(tkn, tkn, KEYWORD)
	} else if tkn[0] >= '0' && tkn[0] <= '9' {
		process(tkn, tkn, INT_CONST)
	} else if tkn[0] == '"' {
		process(tkn, tkn, STRING_CONST)
	} else if tkn == "(" {
		process(tkn, tkn, SYMBOL)
		next = CompileExpression(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
	} else if tkn == "-" || tkn == "~" {
		process(tkn, tkn, SYMBOL)
		next = CompileTerm(content, next)
	} else {
		// need to look ahead to tell 'varName' | 'varName[expression]' | 'subroutineCall'
		tkn2, next, _ = Advance(content, next)
		if tkn2 == "[" {
			next = CompileVarName(content, current)
			tkn, next, _ = Advance(content, next)
			process(tkn, "[", SYMBOL)
			next = CompileExpression(content, next)
			tkn, next, _ = Advance(content, next)
			process(tkn, "]", SYMBOL)
		} else if tkn2 == "." || tkn2 == "(" {
			next = CompileSubroutineCall(content, current)
		} else {
			next = CompileVarName(content, current)
		}
	}
	fmt.Println("</term>")
	return next
}

func CompileSubroutineCall(content []byte, next int) int {
	var tkn Token

	fmt.Println("<subroutineCall>")
	tkn, next, _ = Advance(content, next)
	process(tkn, tkn, IDENTIFIER)
	tkn, next, _ = Advance(content, next)
	if tkn == "." {
		process(tkn, tkn, SYMBOL)
		tkn, next, _ = Advance(content, next)
		process(tkn, tkn, IDENTIFIER)
		tkn, next, _ = Advance(content, next)
		process(tkn, "(", SYMBOL)
		next, _ = CompileExpressionList(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
	} else if tkn == "(" {
		process(tkn, tkn, SYMBOL)
		next, _ = CompileExpressionList(content, next)
		tkn, next, _ = Advance(content, next)
		process(tkn, ")", SYMBOL)
	} else {
		panic(fmt.Sprintf("unterminated subroutineCall at index %d. expected |.| or |(|", next))
	}
	fmt.Println("</subroutineCall>")
	return next
}

func CompileExpressionList(content []byte, next int) (int, int) {
	var tkn Token
	var current, nExprs int

	current = next
	tkn, next, _ = Advance(content, next)
	if tkn == ")" {
		// expressionList is empty
		return current, 0
	}
	fmt.Println("<expressionList>")
	nExprs++
	next = CompileExpression(content, current)
	current = next
	tkn, next, _ = Advance(content, next)
	for tkn == "," {
		nExprs++
		process(tkn, tkn, SYMBOL)
		next = CompileExpression(content, next)
		current = next
		tkn, next, _ = Advance(content, next)
	}
	fmt.Println("</expressionList>")
	return current, nExprs
}

func CompileOp(content []byte, next int) int {
	tkn, next, _ := Advance(content, next)
	process(tkn, tkn, SYMBOL)
	return next
}

func main() {
	content := `
    class X {
        field Y y1;
        static int i1, i2 ;

        constructor X new() { let x = z[11]; return this; }
    function int foo(int j, Arr x) {
        var int b;
        let b = x[4+j];
        let c = x.foo(b);
        return 0;
    }
    method Arr bar() { return null; }
    }
    `
	_ = CompileClass([]byte(content), 0)
}
