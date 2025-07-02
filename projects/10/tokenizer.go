package main

import (
	"fmt"
	"io"
	//"strings"
	"strconv"
	"unicode"
)

type TokenT int
type KeywordT int

type Token = string

const (
	KEYWORD TokenT = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

const (
	CLASS KeywordT = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

func isSymbol(c byte) bool {
	return c == '{' || c == '}' || c == '(' || c == ')' || c == '[' || c == ']' || c == '.' ||
		c == ',' || c == ';' || c == '+' || c == '*' || c == '/' || c == '&' || c == '|' || c == '<' || c == '>' || c == '=' || c == '~'
}

var keywordsMap = map[string]KeywordT{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN}

func isIntConst(s string) (int, error) {
	return strconv.Atoi(s)
}

func Advance(content []byte, start int) (Token, int, error) {
	// Generate the next token from the 'content', starting the tokenize from index 'start'.
	// If no error occurs, returns the Token found, and the index from which starting next.
	if start >= len(content) {
		return "", -1, io.EOF
	}
	end := start
	var c, c2 byte
	for {
		// Skip whitespaces
		if end >= len(content) {
			return "", -1, io.EOF
		}
		c = content[end]
		end++
		if unicode.IsSpace(int32(c)) {
			continue
		}
		if c == '/' {
			// Look ahead for // or /*
			// If found, then this is a comment, hence continue loop
			if end >= len(content) {
				panic(fmt.Sprintf("end of file after '/' at index %d", end))
			}
			c2 = content[end]
			end++
			if c2 == '/' {
				// Skip until '\n'
				for {
					if end >= len(content) {
						// Comment line is also the last line of the file. This is valid.
						break
					}
					if content[end] == '\n' || content[end] == '\r' {
						break
					}
					end++
				}
				continue
			} else if c2 == '*' {
				// Skip until "*/"
				for {
					if end >= len(content)-1 { // we need two chars
						panic(fmt.Sprintf("unterminated comment at index %d", start))
					}
					if content[end] == '*' && content[end+1] == '/' {
						end += 2
						break
					}
					end++
				}
				continue
			} else {
				// Not a comment.
				// Hence '/' is a Symbol on its own
				return "/", end - 1, nil
			}
		} else {
			break
		}
	}
	if isSymbol(c) {
		return string(c), end, nil
	}
	// Else build and return token.
	s1 := end - 1 // start of this token
	for {
		if end >= len(content) {
			panic(fmt.Sprintf("unterminated token at index %d", s1))
		}
		// skip until space or Symbol occurs
		if content[end] == ' ' || isSymbol(content[end]) {
			break
		}
		end++
	}
	tkn := string(content[s1:end])
	return tkn, end, nil
}

func TokenType() TokenT {
	return KEYWORD
}

func Keyword() KeywordT {
	return CLASS
}
