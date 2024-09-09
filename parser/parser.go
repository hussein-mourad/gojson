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
	return p
}

func (p *Parser) Parse() interface{} {
	// document := ast.NewDocument()
	p.eat()

	// JSON must be an array or an object
	if p.currentToken.Type != lexer.LBRACE && p.currentToken.Type != lexer.LBRACKET {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.currentToken.Value, p.currentToken.Line, p.currentToken.Column)
	}
	value := p.parseValue()
	if value == nil {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.currentToken.Value, p.currentToken.Line, p.currentToken.Column)
	}
	return value
}

func (p *Parser) parseValue() interface{} {
	switch p.at().Type {
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.STRING:
		value := p.at().Value
		p.eat()
		return value
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
	// for p.at().Type != lexer.RBRACE {
	// 	if p.isEOF() {
	// 		// TODO: Print Error
	// 	}
	// 	key := p.eat().Value
	//
	// 	if p.currentToken.Type != lexer.COMMA {
	// 		logger.Fatalf("Error: expected : at line: %v column: %v", p.currentToken.Line, p.currentToken.Column)
	// 	}
	// 	p.eat()
	//
	// 	value := p.parseValue()
	//
	// 	obj[key] = value
	//
	// 	if p.currentToken.Type == lexer.COMMA {
	// 		p.eat()
	// 		if p.currentToken.Type == lexer.RBRACE {
	// 			logger.Fatalf("error: unexpected , at line: %v column: %v", p.currentToken.Line, p.currentToken.Column)
	// 		}
	// 	}
	// }
	// p.eat() // Skip closing brace
	// return obj
	return obj
}

func (p *Parser) parseArray() interface{} {
	// TODO: Fix Parsing array eof issues, check pass2.json test
	var arr []interface{}
	p.eat() // skip opening bracket

	for p.currentToken.Type != lexer.RBRACKET {
		value := p.parseValue()
		arr = append(arr, value)

		fmt.Printf("p.curToken: %v\n", p.currentToken.String())
		p.eat()

		if p.currentToken.Type != lexer.COMMA && p.currentToken.Type != lexer.RBRACKET && p.currentToken.Type != lexer.EOF {
			logger.Fatalf("error: unexpected %v at line: %v column: %v", p.currentToken.Value, p.currentToken.Line, p.currentToken.Column)
		}

		// Handle extra comma before the end of the array
		if p.currentToken.Type == lexer.COMMA {
			// Move to the next token after the comma
			p.eat()

			// After a comma, there should be either a value or the end of the array
			if p.currentToken.Type == lexer.RBRACKET {
				logger.Fatalf("error: unexpected %v after comma at line: %v column: %v", p.currentToken.Value, p.currentToken.Line, p.currentToken.Column)
			}
		}

		// If the token is lexer.EOF, handle it appropriately
		if p.currentToken.Type == lexer.EOF {
			// End of file without closing bracket
			logger.Fatalf("error: unexpected end of file while parsing array at line: %v column: %v", p.currentToken.Line, p.currentToken.Column)
			return arr // Or handle the error based on your needs
		}
	}

	p.eat() // skip closing bracket
	return arr
}

func (p *Parser) parseNumber() interface{} {
	s := p.currentToken.Value

	// Try to parse as int
	if intVal, err := strconv.Atoi(s); err == nil {
		if len(strconv.Itoa(intVal)) != len(s) {
			// Check if number doesn't have leading zeros
			logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.currentToken.Line, p.currentToken.Column)
		}
		p.eat()
		return intVal
	}
	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		p.eat()
		return floatVal
	}

	logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.currentToken.Line, p.currentToken.Column)
	return nil
}

func (p *Parser) parseBoolean() bool {
	return p.eat().Value == "true"
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

func (p *Parser) expect(Type lexer.TokenType, err string) {
	if p.at().Type != Type {
		log.Fatalf("Error: %v on line %v, column %v", err, p.at().Line, p.at().Column)
	}
}
