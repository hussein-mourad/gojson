package main

import (
	"fmt"
	"os"
	"unicode"
)

const (
	LEFTBRACE    = "LEFTBRACE"
	RIGHTBRACE   = "RIGHTBRACE"
	LEFTBRACKET  = "LEFTBRACKET"
	RIGHTBRACKET = "RIGHTBRACKET"
	COMMA        = "COMMA"
	COLON        = "COLON"
	STRING       = "STRING"
	NUMBER       = "NUMBER"
	TRUE         = "TRUE"
	FALSE        = "FALSE"
	NULL         = "NULL"
	EOF          = "EOF"
	UNKOWN       = "UNKOWN"
)

type Token struct {
	Type  string
	Value string
}

func NewToken(Type string, Value string) *Token {
	return &Token{Type, Value}
}

func (t *Token) String() string {
	return fmt.Sprintf("Type: %v\tValue: %v", t.Type, t.Value)
}

type Lexer struct {
	input       string // code
	pos         int    // position
	currentRune rune   // current character
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readRune()
	return l
}

func (l *Lexer) makeToken(Type string) *Token {
	value := string(l.currentRune)
	l.readRune()
	return NewToken(Type, value)
}

func (l *Lexer) readRune() {
	if l.pos >= len(l.input) {
		l.currentRune = 0 // EOF
	} else {
		l.currentRune = rune(l.input[l.pos])
	}
	l.pos++
}

func (l *Lexer) readString() *Token {
	var str string
	l.readRune() // skip opening qoute
	for l.currentRune != '"' && l.currentRune != 0 {
		str += string(l.currentRune)
		l.readRune()
	}
	l.readRune() // skip closing quote
	return NewToken(STRING, str)
}

func (l *Lexer) readNumber() *Token {
	var str string
	for unicode.IsDigit(l.currentRune) || l.currentRune == '-' || l.currentRune == '.' {
		str += string(l.currentRune)
		l.readRune()
	}
	return NewToken(STRING, str)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.currentRune) {
		l.readRune()
	}
}

func (l *Lexer) NextToken() *Token {
	l.skipWhitespace()
	switch l.currentRune {
	case '{':
		return l.makeToken(LEFTBRACE)
	case '}':
		return l.makeToken(RIGHTBRACE)
	case '[':
		return l.makeToken(LEFTBRACKET)
	case ']':
		return l.makeToken(RIGHTBRACKET)
	case ':':
		return l.makeToken(COLON)
	case ',':
		return l.makeToken(COMMA)
	case '"':
		return l.readString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return l.readNumber()
	}

	if l.currentRune == 0 {
		return NewToken(EOF, "")
	}

	return l.makeToken(UNKOWN)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
	}

	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	lexer := NewLexer(string(data))
	for {
		t := lexer.NextToken()
		fmt.Printf("%v\n", t)
		if t.Type == EOF {
			break
		}
	}
}
