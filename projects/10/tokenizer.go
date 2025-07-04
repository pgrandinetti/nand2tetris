package main

import (
	"fmt"
	"io"
	//"strings"
	"regexp"
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
		c == ',' || c == ';' || c == '+' || c == '-' || c == '*' || c == '/' || c == '&' || c == '|' || c == '<' || c == '>' || c == '=' || c == '~'
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

func isKeyword(s string) bool {
	_, ok := keywordsMap[s]
	return ok
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
		// Skip whitespaces and comments
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
	if c == '"' {
		// skip until next "
		for {
			if end >= len(content) {
				panic(fmt.Sprintf("unterminated string at index %d", s1))
			}
			if content[end] == '\n' {
				panic(fmt.Sprintf("invalid string constant at index %d", s1))
			}
			if content[end] == '"' {
				break
			}
			end++
		}
		return string(content[s1 : end+1]), end + 1, nil
	}
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

func TokenType(tkn Token) TokenT {
	if len(tkn) == 1 && isSymbol(tkn[0]) {
		return SYMBOL
	}
	if isKeyword(tkn) {
		return KEYWORD
	}
	if _, err := isIntConst(tkn); err == nil {
		return INT_CONST
	}
	if tkn[0] == '"' && tkn[len(tkn)-1] == '"' {
		return STRING_CONST
	}
	pattern := `^[a-zA-Z_][a-zA-Z0-9_]*$`
	match, _ := regexp.MatchString(pattern, tkn)
	if match {
		return IDENTIFIER
	}
	panic(fmt.Sprintf("token not valid: %s", tkn))
}

func Keyword(tkn Token) KeywordT {
	val, ok := keywordsMap[tkn]
	if !ok {
		panic(fmt.Sprintf("token is not a keyword: %s", tkn))
	}
	return val
}
