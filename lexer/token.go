package lexer

import "fmt"

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
