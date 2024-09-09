package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hussein-mourad/go-json-parser/ast"
	"github.com/hussein-mourad/go-json-parser/lexer"
)

var logger = log.New(os.Stderr, "", 0)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken *lexer.Token // current Token
	ast          *ast.Document
	data         interface{}
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, ast: ast.NewDocument()}
	p.eat()
	return p
}

func (p *Parser) Parse() interface{} {
	p.notExpect(lexer.EOF, "Unexpected end of file")
	p.expectIn("Json payload should be an object or a string", lexer.LBRACE, lexer.LBRACKET)
	p.data, p.ast.Body = p.parseValue()
	p.expect(lexer.EOF, "Unexpected character")
	return p.data
}

func (p *Parser) parseValue() (interface{}, ast.Stmt) {
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
		return nil, nil
	}
	return nil, nil
}

func (p *Parser) parseObject() (map[string]interface{}, *ast.Object) {
	obj := make(map[string]interface{})
	objStmt := ast.NewObject()
	p.eat() // Skip opening brace
	for p.at().Type != lexer.RBRACE {
		p.notExpect(lexer.EOF, "Unexpected end of file")
		key := p.eat().Value
		p.expect(lexer.COLON, "Expected :")
		p.eat()
		value, stmt := p.parseValue()
		obj[key] = value
		property := ast.NewProperty()
		property.Key.Value = key
		property.Value = stmt
		objStmt.Members = append(objStmt.Members, property)
		p.notExpect(lexer.EOF, "Unexpected end of file")
		p.expectIn("Unexpected character", lexer.COMMA, lexer.RBRACE)
		// If there is a comma then there should be new values
		if p.currentToken.Type == lexer.COMMA {
			p.eat()
			p.expect(lexer.STRING, "Unexpected character")
		}
	}
	p.eat() // Skip closing brace
	return obj, objStmt
}

func (p *Parser) parseArray() ([]interface{}, *ast.Array) {
	var arr []interface{}
	arrStmt := ast.NewArray()
	p.eat() // skip opening bracket
	for p.currentToken.Type != lexer.RBRACKET {
		p.notExpect(lexer.EOF, "Unexpected end of file")
		value, stmt := p.parseValue()
		arr = append(arr, value)
		arrStmt.Elements = append(arrStmt.Elements, stmt)
		p.notExpect(lexer.EOF, "Unexpected end of file")
		p.expectIn(fmt.Sprintf("Unexpected %v", p.at().Value), lexer.COMMA, lexer.RBRACKET)
		if p.currentToken.Type == lexer.COMMA {
			p.eat()
			p.notExpect(lexer.RBRACKET, "Unexpected character")
		}
	}
	p.eat() // skip closing bracket
	return arr, arrStmt
}

func (p *Parser) parseString() (interface{}, *ast.StringLiteral) {
	value := p.eat().Value
	return value, ast.NewStringLiteral(value)
}

func (p *Parser) parseBoolean() (interface{}, *ast.BooleanLiteral) {
	value := p.eat().Value == "true"
	return value, ast.NewBooleanLiteral(value)
}

func (p *Parser) parseNumber() (interface{}, *ast.NumberLiteral) {
	value := p.eat().Value

	intVal, err := p.parserInt(value)
	if err == nil {
		return intVal, ast.NewNumberLiteral(intVal)
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return floatVal, ast.NewNumberLiteral(floatVal)
	}

	complexVal, err := strconv.ParseComplex(value, 128)
	if err == nil {
		return complexVal, ast.NewNumberLiteral(complexVal)
	}
	p.error("Can't parse number")
	return nil, nil
}

func (p *Parser) parserInt(value string) (int64, error) {
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		// Check if ints have leading zeros
		if len(strconv.Itoa(int(intVal))) != len(value) {
			p.error("Numbers cannot have leading zeros")
		}
		return intVal, nil
	}
	return 0, err
}

func (p *Parser) GetAST() *ast.Document {
	return p.ast
}

func (p *Parser) isEOF() bool {
	return p.currentToken.Type == lexer.EOF
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

func (p *Parser) expectIn(err string, Types ...lexer.TokenType) {
	found := false
	for _, t := range Types {
		if p.at().Type == t {
			found = true
			break
		}
	}
	if !found {
		p.error(err)
	}
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
	logger.Fatalf("Error: %v on line %v, column %v", msg, p.at().Line, p.at().Column)
}
