package lexer

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

var EOFCHAR rune = 0

var logger = log.New(os.Stderr, "", 0)

var CharTokens = map[rune]TokenType{
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	':': COLON,
	',': COMMA,
}

var KeywordTokens = map[string]TokenType{
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

func (l *Lexer) makeToken(Type TokenType, Value string) *Token {
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
	if char == EOFCHAR {
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
		if l.at() == EOFCHAR {
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
			} else {
				l.error("Unexpected characater")
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
	for unicode.IsDigit(l.at()) || l.IsOneOfMany('-', '.', 'e', 'E', '+') {
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
	char := l.at()
	if char == '\n' {
		l.column = 1
		l.line++
	}
	l.current = string(char)
	return char
}

func (l *Lexer) advanceN(n int) rune {
	if l.index+n < len(l.input) {
		l.index += n
		l.column += n
		l.current = string(l.at())
		return l.at()
	}
	return EOFCHAR
}

func (l *Lexer) at() rune {
	if l.index < len(l.input) {
		return rune(l.input[l.index])
	}
	return EOFCHAR
}

func (l *Lexer) IsOneOfMany(expected ...rune) bool {
	// Is the current character in the list of characters
	for _, e := range expected {
		if l.at() == e {
			return true
		}
	}
	return false
}

func (l *Lexer) errorEOF() {
	l.error("Unexpected end of file")
}

func (l *Lexer) errorUnexpectedChar(char rune) {
	l.error(fmt.Sprintf("Unexpected character %c", char))
}

func (l *Lexer) error(msg string) {
	logger.Fatalf("Error: %v at line: %v, column: %v\n", msg, l.line, l.column)
}
