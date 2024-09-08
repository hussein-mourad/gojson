package main

import (
	"fmt"
	"os"
	"unicode"
)

const (
	LBRACE   = iota // {
	RBRACE          // }
	LBRACKET        // [
	RBRACKET        // ]
	COLON           // :
	COMMA           // ,

	STRING
	NUMBER
	BOOLEAN
	NULL

	EOF
)

var CharTokens = map[rune]int{
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	':': COLON,
	',': COMMA,
}

var KeywordTokens = map[string]int{
	"true":  BOOLEAN,
	"false": BOOLEAN,
	"null":  NULL,
}

var EscapeChar = map[rune]rune{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

var NumberChars = map[rune]rune{
	'0': '0',
	'1': '1',
	'2': '2',
	'3': '3',
	'4': '4',
	'5': '5',
	'6': '6',
	'7': '7',
	'8': '8',
	'9': '9',
	'-': '-',
}

type Token struct {
	Type   int
	Value  string
	Line   int
	Column int
	Index  int
}

func newToken(Type int, Value string, Line int, Column int, Index int) *Token {
	return &Token{Type, Value, Line, Column, Index}
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %v, Value: %v, Line: %v, Column: %v, Index: %v}", t.TypeString(), t.Value, t.Line, t.Column, t.Index)
}

func (t *Token) TypeString() string {
	var str string
	switch t.Type {
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case STRING:
		return "String"
	case NUMBER:
		return "Number"
	case BOOLEAN:
		return "Boolean"
	case NULL:
		return "Null"
	case EOF:
		return "EOF"
	}
	return str
}

type Lexer struct {
	input   string // code
	index   int    // index in code
	line    int    // current line
	column  int    // current column in line
	current string // current character (used for debugging)
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 1}
	l.current = string(l.at())
	return l
}

func (l *Lexer) makeToken(Type int, Value string) *Token {
	return newToken(Type, Value, l.line, l.column, l.index)
}

func (l *Lexer) NextToken() *Token {
	l.skipWhitespaces()
	// Order matters
	parsers := []func() *Token{l.readChar, l.readNumbers, l.readKeywords, l.readString}
	var matchedToken *Token
	for _, parser := range parsers {
		matchedToken = parser()
		if matchedToken != nil {
			break
		}
	}
	return matchedToken
}

func (l *Lexer) readChar() *Token {
	char := l.at()
	if char == 0 {
		return l.makeToken(EOF, "")
	}
	Type, exists := CharTokens[char]
	if !exists {
		return nil
	}
	token := l.makeToken(Type, string(char))
	l.advance()
	return token
}

func (l *Lexer) readString() *Token {
	str := ""
	found := false
	if l.at() != '"' {
		l.errorUnexpectedChar(l.at())
	}
	l.advance() // skip opening quote
	found = true
	for l.at() != '"' {
		if l.at() == 0 {
			l.errorEOF()
		}
		s := l.at()
		if l.at() == '\\' {
			l.advance()
			c, exists := EscapeChar[l.at()]
			if exists {
				s = c
			} else if l.at() == 'u' { // unicode
				// add \u to the string
				str += string('\\')
				s = 'u'
			}
		}
		str += string(s)
		l.advance()
	}

	if l.at() != '"' {
		l.errorUnexpectedChar(l.at())
	}
	l.advance() // skip closing quote
	if !found {
		return nil
	}
	return l.makeToken(STRING, str)
}

func (l *Lexer) readKeywords() *Token {
	var token *Token
	for keyword, Type := range KeywordTokens {
		kwLen := len(keyword)
		if l.index+kwLen >= len(l.input) {
			continue
		}
		if l.input[l.index:l.index+kwLen] == keyword {
			token = l.makeToken(Type, keyword)
			l.advanceN(kwLen)
		}
	}
	return token
}

func (l *Lexer) readNumbers() *Token {
	if !unicode.IsDigit(l.at()) && l.at() != '-' {
		return nil
	}
	var num string
	for unicode.IsDigit(l.at()) || IsOneOfMany(l.at(), '-', '.', 'e', 'E', '+') {
		num += string(l.at())
		l.advance()
	}
	return l.makeToken(NUMBER, num)
}

func (l *Lexer) skipWhitespaces() {
	for unicode.IsSpace(l.at()) {
		l.advance()
	}
}

func (l *Lexer) advance() rune {
	l.index++
	l.column++
	if l.index < len(l.input) {
		r := rune(l.input[l.index])
		l.current = string(r)
		char := l.at()
		if char == 0 {
			return 0
		}
		if char == '\n' {
			l.column = 1
			l.line++
		}
		return r
	}
	return 0
}

func (l *Lexer) advanceN(n int) rune {
	if l.index+n < len(l.input) {
		l.index += n
		l.column += n
		if l.index < len(l.input) {
			l.current = string(l.input[l.index])
		}
		return l.at()
	}
	return 0
}

func (l *Lexer) at() rune {
	if l.index < len(l.input) {
		return rune(l.input[l.index])
	}
	return 0
}

func (l *Lexer) errorEOF() {
	fmt.Printf("Error: Unexpected end of file at line: %v, column: %v\n", l.line, l.column)
	os.Exit(1)
}

func (l *Lexer) errorUnexpectedChar(char rune) {
	fmt.Printf("Error: Unexpected character %c at line: %v, column: %v\n", char, l.line, l.column)
	os.Exit(1)
}

func (l *Lexer) error(msg string) {
	fmt.Printf("Error: %v at line: %v, column: %v\n", msg, l.line, l.column)
	os.Exit(1)
}

func IsOneOfMany(value rune, expected ...rune) bool {
	for _, e := range expected {
		if value == e {
			return true
		}
	}
	return false
}
