package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"unicode"
)

const (
	LEFTBRACE    = "{"
	RIGHTBRACE   = "}"
	LEFTBRACKET  = "["
	RIGHTBRACKET = "]"
	COMMA        = ","
	COLON        = ";"
	STRING       = "STRING"
	NUMBER       = "NUMBER"
	TRUE         = "TRUE"
	FALSE        = "FALSE"
	NULL         = "NULL"
	EOF          = "EOF"
)

type Token struct {
	Type  string
	Value string
}

type Lexer struct {
	input       string // code
	pos         int    // position
	currentRune rune   // current character
}

func NewToken(Type string, Value string) *Token {
	return &Token{Type, Value}
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) readRune() {
	if l.pos >= len(l.input) {
		l.currentRune = 0 // EOF
	} else {
		l.currentRune = rune(l.input[l.pos])
	}
	l.pos++
}

func tokenize(line string) ([]Token, error) {
	var tokens []Token
	for _, r := range line {
		if r == '{' {
			tokens = append(tokens, *NewToken(LEFTBRACE, string(r)))
		}
		if r == '}' {
			tokens = append(tokens, *NewToken(RIGHTBRACE, string(r)))
			return tokens, nil
		}
	}
	return nil, errors.New("invalid line")
}

func skipWhitespaces(reader *bufio.Reader) rune {
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if !unicode.IsSpace(r) {
			return r
		}
	}
	return 0
}

func readString(reader *bufio.Reader) string {
	var str string
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(r) {
			continue
		}
		if r == '"' {
			break
		}
		str += string(r)
	}
	return str
}

func main() {
	// panic("Test panic")
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s [file.json]\n", os.Args[0])
	}

	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(file)
	for {
		var str string
		r := skipWhitespaces(reader)
		if r == 0 {
			break
		}

		switch r {
		case '"':
			str = readString(reader)
		default:
			str = string(r)
		}

		fmt.Printf("token: %v\n", str)
	}

	// r, _, _ := reader.ReadRune()
	// r, _, _ = reader.ReadRune()
	// fmt.Printf("r: %c\n", r)

	// for {
	// 	r, _, err := reader.ReadRune()
	// 	if err != nil {
	// 		break
	// 	}
	// 	switch r {
	// 	case '{':
	// 		fmt.Println("Left braces")
	// 	case '}':
	// 		fmt.Println("Right braces")
	// 	case '"':
	// 		fmt.Println("Right braces")
	// 	case ':':
	// 		fmt.Println("Right braces")
	// 	}
	// }

	// fmt.Println()
	// s := bufio.NewScanner(file)
	// s.Split(bufio.ScanRunes)
	// for s.Scan() {
	// 	r := s.Text()
	// 	fmt.Printf("r: %v\n", r)
	// tokens, err := tokenize(line)
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// 	os.Exit(1)
	// }
	// for _, token := range tokens {
	// 	fmt.Printf("Type: %v\tValue: %v\n", token.Type, token.Value)
	// }
	// }
}
