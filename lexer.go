package main

import (
	"fmt"
	"unicode"
)

const (
	LBRACE   = "LBRACE"
	RBRACE   = "RBRACE"
	LBRACKET = "LBRACKET"
	RBRACKET = "RBRACKET"
	COMMA    = "COMMA"
	COLON    = "COLON"
	STRING   = "STRING"
	NUMBER   = "NUMBER"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	EOF      = "EOF"
	UNKOWN   = "UNKOWN"
)

type Token struct {
	Type   string
	Value  string
	Line   int
	Column int
}

func NewToken(Type string, Value string, Line int, Column int) *Token {
	return &Token{Type, Value, Line, Column}
}

func (t *Token) String() string {
	return fmt.Sprintf("Type: %v\tValue: %v\tLine: %v\tColumn: %v", t.Type, t.Value, t.Line, t.Column)
}

type Lexer struct {
	input       string // code
	pos         int    // position
	currentRune rune   // current character
	line        int    // current line
	column      int    // current column
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.readRune()
	return l
}

func (l *Lexer) makeToken(Type string) *Token {
	value, line, column := string(l.currentRune), l.line, l.column
	l.readRune()
	return NewToken(Type, value, line, column)
}

func (l *Lexer) readRune() {
	if l.pos >= len(l.input) {
		l.currentRune = 0 // EOF
	} else {
		l.currentRune = rune(l.input[l.pos])
	}

	// fmt.Printf("ch: %q\tline: %v\n", l.currentRune, l.line)

	if l.currentRune == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.pos++
}

func (l *Lexer) readString() *Token {
	var str string
	line, column := l.line, l.column
	l.readRune() // skip opening qoute
	for l.currentRune != '"' && l.currentRune != 0 {
		str += string(l.currentRune)
		l.readRune()
	}
	l.readRune() // skip closing quote
	return NewToken(STRING, str, line, column)
}

func (l *Lexer) readNumber() *Token {
	var str string
	line, column := l.line, l.column
	for unicode.IsDigit(l.currentRune) || l.currentRune == '-' || l.currentRune == '.' {
		str += string(l.currentRune)
		l.readRune()
	}
	return NewToken(NUMBER, str, line, column)
}

func (l *Lexer) readStringWithoutQuotes() *Token {
	var str string
	line, column := l.line, l.column
	tokenType := UNKOWN
	// '\n' is for handling last record that may not have a comma
	for l.currentRune != ',' && l.currentRune != '\n' && l.currentRune != 0 {
		str += string(l.currentRune)
		l.readRune()
	}
	switch str {
	case "true":
		tokenType = TRUE
	case "false":
		tokenType = FALSE
	case "null":
		tokenType = NULL
	}
	return NewToken(tokenType, str, line, column)
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
		return l.makeToken(LBRACE)
	case '}':
		return l.makeToken(RBRACE)
	case '[':
		return l.makeToken(LBRACKET)
	case ']':
		return l.makeToken(RBRACKET)
	case ':':
		return l.makeToken(COLON)
	case ',':
		return l.makeToken(COMMA)
	case '"':
		return l.readString()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		return l.readNumber()

	}
	switch {
	case unicode.IsLetter(l.currentRune):
		return l.readStringWithoutQuotes()
	}

	if l.currentRune == 0 {
		return NewToken(EOF, "", l.line, l.column)
	}

	return l.makeToken(UNKOWN)
}
