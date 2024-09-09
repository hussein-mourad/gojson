package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hussein-mourad/go-json-parser/ast"
	"github.com/hussein-mourad/go-json-parser/lexer"
)

var logger = log.New(os.Stdout, "", 0)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *lexer.Token // current Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.eat()
	return p
}

func (p *Parser) Parse() *ast.Document {
	document := ast.NewDocument()
	document.Body = p.parseValue()
	p.expect(lexer.EOF, "Unexpected character")
	return document
}

func (p *Parser) parseValue() ast.Stmt {
	switch p.at().Type {
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.STRING:
		return p.parseString()
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.BOOLEAN:
		return p.parseBoolean()
	case lexer.NULL:
		p.eat()
		return nil
	}
	return nil
}

func (p *Parser) parseObject() *ast.Object {
	obj := ast.NewObject()
	p.eat() // Skip opening brace
	for p.at().Type != lexer.RBRACE {
		p.notExpect(lexer.EOF, "Unexpected end of file")
		property := ast.NewProperty()
		property.Key.Value = p.eat().Value
		p.expect(lexer.COLON, "Expected :")
		p.eat()
		property.Value = property.Value
		obj.Members = append(obj.Members, property)
		// If there is a comma then there should be new values
		if p.at().Type == lexer.COMMA {
			p.eat()
			p.notExpect(lexer.RBRACE, "Unexpected }")
			p.notExpect(lexer.RBRACKET, "Unexpected ]")
		}
	}
	p.eat() // Skip closing brace
	return obj
}

func (p *Parser) parseArray() *ast.Array {
	arr := ast.NewArray()
	p.eat() // skip opening bracket
	for p.at().Type != lexer.RBRACKET {
		p.notExpect(lexer.EOF, "Unexpected end of file")
		arr.Elements = append(arr.Elements, p.parseValue())
		if p.at().Type == lexer.COMMA {
			p.eat()
			p.notExpect(lexer.RBRACE, "Unexpected }")
			p.notExpect(lexer.RBRACKET, "Unexpected ]")
		}
	}
	p.eat() // skip closing bracket
	return arr
}

func (p *Parser) parseString() *ast.StringLiteral {
	return ast.NewStringLiteral(p.eat().Value)
}

func (p *Parser) parseBoolean() *ast.BooleanLiteral {
	return ast.NewBooleanLiteral(p.eat().Value == "true")
}

func (p *Parser) parseNumber() *ast.NumberLiteral {
	value := p.eat().Value
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return ast.NewNumberLiteral(intVal)
	}
	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return ast.NewNumberLiteral(floatVal)
	}

	complexVal, err := strconv.ParseComplex(value, 128)
	if err == nil {
		return ast.NewNumberLiteral(complexVal)
	}
	p.error("Can't parse number")
	return nil
}

func (p *Parser) isEOF() bool {
	return p.at().Type == lexer.EOF
}

func (p *Parser) at() *lexer.Token {
	return p.currentToken
}

func (p *Parser) eat() *lexer.Token {
	// Returns the previous token and then advances
	prev := p.currentToken
	p.currentToken = p.lexer.NextToken()
	fmt.Println(p.currentToken) // FIXME: debugging only
	return prev
}

func (p *Parser) expect(Type lexer.TokenType, err string) {
	if p.at().Type != Type {
		p.error(err)
	}
}

func (p *Parser) notExpect(Type lexer.TokenType, err string) {
	if p.at().Type == Type {
		p.error(err)
	}
}

func (p *Parser) error(msg string) {
	log.Fatalf("Error: %v on line %v, column %v", msg, p.at().Line, p.at().Column)
}
