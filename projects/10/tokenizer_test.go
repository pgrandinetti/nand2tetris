package main

import (
	"io"
	"testing"
)

func TestAdvance(t *testing.T) {
	inpt := `if (bar<   10)  {
        // comment
        let bla = 5;
        /* another
     comment*/
        x /y  // to test the single '/'
     }
     let s = "  a string"   ;
     //comment to end file
    `
	inptb := []byte(inpt)
	start := 0
	exp := []map[string]any{
		{"tkn": "if", "next": 2, "err": nil},
		{"tkn": "(", "next": 4, "err": nil},
		{"tkn": "bar", "next": 7, "err": nil},
		{"tkn": "<", "next": 8, "err": nil},
		{"tkn": "10", "next": 13, "err": nil},
		{"tkn": ")", "next": 14, "err": nil},
		{"tkn": "{", "next": 17, "err": nil},
		{"tkn": "let", "next": 48, "err": nil},
		{"tkn": "bla", "next": 52, "err": nil},
		{"tkn": "=", "next": 54, "err": nil},
		{"tkn": "5", "next": 56, "err": nil},
		{"tkn": ";", "next": 57, "err": nil},
		{"tkn": "x", "next": 101, "err": nil},
		{"tkn": "/", "next": 103, "err": nil},
		{"tkn": "y", "next": 104, "err": nil},
		{"tkn": "}", "next": 138, "err": nil},
		{"tkn": "let", "next": 147, "err": nil},
		{"tkn": "s", "next": 149, "err": nil},
		{"tkn": "=", "next": 151, "err": nil},
		{"tkn": "\"  a string\"", "next": 164, "err": nil},
		{"tkn": ";", "next": 168, "err": nil},
		{"tkn": "", "next": -1, "err": io.EOF}}
	for i := 0; i < len(exp); i++ {
		tkn, next, err := Advance(inptb, start)
		if tkn != exp[i]["tkn"] {
			t.Errorf("expected '%s' computed |%s|", exp[i]["tkn"], tkn)
		}
		if next != exp[i]["next"] {
			t.Errorf("expected next=%d computed next=%d", exp[i]["next"], next)
		}
		if err != exp[i]["err"] {
			t.Errorf("expected error %v computed error %v", exp[i]["err"], err)
		}
		start = next
	}
}

func TestUnterminatedString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Advance did not panic")
		}
	}()
	inpt := `
        let x = "this is not finishe
        d"
        `
	inptb := []byte(inpt)
	var err error = nil
	start := 0
	for err == nil {
		_, start, err = Advance(inptb, start)
	}
}

func TestUnterminatedComent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Advance did not panic")
		}
	}()
	inpt := `
        let x = /*this is not finishe
        d"
        `
	inptb := []byte(inpt)
	var err error = nil
	start := 0
	for err == nil {
		_, start, err = Advance(inptb, start)
	}
}

func TestTokenType(t *testing.T) {
	inpt := "foo"
	exp := IDENTIFIER
	res := TokenType(inpt)
	if res != exp {
		t.Errorf("token %s expected type %v computed %v", inpt, exp, res)
	}
}
