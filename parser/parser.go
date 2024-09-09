package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/hussein-mourad/go-json-parser/lexer"
)

var logger = log.New(os.Stdout, "", 0)

type Parser struct {
	lexer    *lexer.Lexer
	curToken lexer.Token // current Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = *p.lexer.NextToken()
	fmt.Println(p.curToken.String()) // Debug
	// if p.curToken.Type == UNKOWN {
	// 	logger.Fatalf("error: unexpected token %v at line %v column %v\n", p.curToken.Value, p.curToken.Line, p.curToken.Column)
	// }
}

func (p *Parser) parseNumber() interface{} {
	s := p.curToken.Value

	// Try to parse as int
	if intVal, err := strconv.Atoi(s); err == nil {
		if len(strconv.Itoa(intVal)) != len(s) {
			// Check if number doesn't have leading zeros
			logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
		return intVal
	}
	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		p.nextToken()
		return floatVal
	}

	logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.curToken.Line, p.curToken.Column)
	return nil
}

func (p *Parser) parse() interface{} {
	// JSON must be an array or an object
	if p.curToken.Type != lexer.LBRACE && p.curToken.Type != lexer.LBRACKET {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.curToken.Line, p.curToken.Column)
	}
	value := p.parseValue()
	if value == nil {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.curToken.Line, p.curToken.Column)
	}
	return value
}

func (p *Parser) parseValue() interface{} {
	switch p.curToken.Type {
	case lexer.STRING:
		value := p.curToken.Value
		p.nextToken()
		return value
	case lexer.NUMBER:
		return p.parseNumber()
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	case lexer.BOOLEAN:
		p.nextToken()
		// TODO: parse false
		return true
	case lexer.NULL:
		p.nextToken()
		return nil
	}
	return nil
}

func (p *Parser) parseObject() interface{} {
	obj := make(map[string]interface{})
	p.nextToken() // Skip opening brace
	for p.curToken.Type != lexer.RBRACE {
		key := p.curToken.Value
		p.nextToken()

		if p.curToken.Value != ":" {
			logger.Fatalf("error: expected : at line: %v column: %v", p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()

		value := p.parseValue()

		obj[key] = value

		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
			if p.curToken.Type == lexer.RBRACE {
				logger.Fatalf("error: unexpected , at line: %v column: %v", p.curToken.Line, p.curToken.Column)
			}
		}
	}
	p.nextToken() // Skip closing brace
	return obj
}

func (p *Parser) parseArray() interface{} {
	// TODO: Fix Parsing array eof issues, check pass2.json test
	var arr []interface{}
	p.nextToken() // skip opening bracket

	for p.curToken.Type != lexer.RBRACKET {
		value := p.parseValue()
		arr = append(arr, value)

		fmt.Printf("p.curToken: %v\n", p.curToken.String())
		p.nextToken()

		if p.curToken.Type != lexer.COMMA && p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
			logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.curToken.Line, p.curToken.Column)
		}

		// Handle extra comma before the end of the array
		if p.curToken.Type == lexer.COMMA {
			// Move to the next token after the comma
			p.nextToken()

			// After a comma, there should be either a value or the end of the array
			if p.curToken.Type == lexer.RBRACKET {
				logger.Fatalf("error: unexpected %v after comma at line: %v column: %v", p.curToken.Value, p.curToken.Line, p.curToken.Column)
			}
		}

		// If the token is lexer.EOF, handle it appropriately
		if p.curToken.Type == lexer.EOF {
			// End of file without closing bracket
			logger.Fatalf("error: unexpected end of file while parsing array at line: %v column: %v", p.curToken.Line, p.curToken.Column)
			return arr // Or handle the error based on your needs
		}
	}

	p.nextToken() // skip closing bracket
	return arr
}
