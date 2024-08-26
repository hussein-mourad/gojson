package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer    *Lexer
	curToken Token // current Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = *p.lexer.NextToken()
	fmt.Println(p.curToken.String()) // Debug
	if p.curToken.Type == UNKOWN {
		logger.Fatalf("error: unexpected token %v at line %v column %v\n", p.curToken.Value, p.lexer.line, p.lexer.column)
	}
}

func (p *Parser) parseNumber() interface{} {
	s := p.curToken.Value

	// Try to parse as int
	if intVal, err := strconv.Atoi(s); err == nil {
		if len(strconv.Itoa(intVal)) != len(s) {
			// Check if number doesn't have leading zeros
			logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.lexer.line, p.lexer.column)
		}
		p.nextToken()
		return intVal
	}
	// Try to parse as float
	if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		p.nextToken()
		return floatVal
	}

	logger.Fatalf("error: invalid number %v at line: %v column: %v", s, p.lexer.line, p.lexer.column)
	return nil
}

func (p *Parser) parse() interface{} {
	// JSON must be an array or an object
	if p.curToken.Type != LBRACE && p.curToken.Type != LBRACKET {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.lexer.line, p.lexer.column)
	}
	value := p.parseValue()
	if value == nil {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.lexer.line, p.lexer.column)
	}
	return value
}

func (p *Parser) parseValue() interface{} {
	switch p.curToken.Type {
	case STRING:
		value := p.curToken.Value
		p.nextToken()
		return value
	case NUMBER:
		return p.parseNumber()
	case LBRACE:
		return p.parseObject()
	case LBRACKET:
		return p.parseArray()
	case TRUE:
		p.nextToken()
		return true
	case FALSE:
		p.nextToken()
		return false
	case NULL:
		p.nextToken()
		return nil
	}
	return nil
}

func (p *Parser) parseObject() interface{} {
	obj := make(map[string]interface{})
	p.nextToken() // Skip opening brace
	for p.curToken.Type != RBRACE {
		key := p.curToken.Value
		p.nextToken()

		if p.curToken.Value != ":" {
			logger.Fatalf("error: expected : at line: %v column: %v", p.lexer.line, p.lexer.column)
		}
		p.nextToken()

		value := p.parseValue()

		obj[key] = value

		if p.curToken.Type == COMMA {
			p.nextToken()
			if p.curToken.Type == RBRACE {
				logger.Fatalf("error: unexpected , at line: %v column: %v", p.lexer.line, p.lexer.column)
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

	for p.curToken.Type != RBRACKET {
		value := p.parseValue()
		arr = append(arr, value)

		fmt.Printf("p.curToken: %v\n", p.curToken.String())
		p.nextToken()

		if p.curToken.Type != COMMA && p.curToken.Type != RBRACKET && p.curToken.Type != EOF {
			logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.lexer.line, p.lexer.column)
		}

		// Handle extra comma before the end of the array
		if p.curToken.Type == COMMA {
			// Move to the next token after the comma
			p.nextToken()

			// After a comma, there should be either a value or the end of the array
			if p.curToken.Type == RBRACKET {
				logger.Fatalf("error: unexpected %v after comma at line: %v column: %v", p.curToken.Value, p.lexer.line, p.lexer.column)
			}
		}

		// If the token is EOF, handle it appropriately
		if p.curToken.Type == EOF {
			// End of file without closing bracket
			logger.Fatalf("error: unexpected end of file while parsing array at line: %v column: %v", p.lexer.line, p.lexer.column)
			return arr // Or handle the error based on your needs
		}
	}

	p.nextToken() // skip closing bracket
	return arr
}
