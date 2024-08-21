package main

import (
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
	if p.curToken.Type == UNKOWN {
		logger.Fatalf("error: unexpected token %v at line %v column %v\n", p.curToken.Value, p.lexer.line, p.lexer.column)
	}
}

func isValidNumber(num string) bool {
	v, err := strconv.Atoi(num)
	if err != nil {
		return false
	}

	if len(strconv.Itoa(v)) != len(num) {
		// Check if number has leading zeros
		return false
	}
	return true
}

func (p *Parser) parse() interface{} {
	value := p.parseValue()
	if value == nil {
		logger.Fatalf("error: unexpected %v at line: %v column: %v", p.curToken.Value, p.lexer.line, p.lexer.column)
	}
	return value
	// if p.curToken.Type == LBRACE {
	// 	return p.parseObject()
	// } else {
	// 	logger.Fatalf("error: expected { at line: %v column: %v", p.lexer.line, p.lexer.column)
	// }
	// return nil
}

func (p *Parser) parseValue() interface{} {
	switch p.curToken.Type {
	case STRING:
		value := p.curToken.Value
		p.nextToken()
		return value
	case NUMBER:
		value := p.curToken.Value
		if !isValidNumber(value) {
			logger.Fatalf("error: invalid number %v at line: %v column: %v", value, p.lexer.line, p.lexer.column)
		}
		p.nextToken()
		return value
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
	var arr []interface{}

	p.nextToken() // skip opening bracket

	for p.curToken.Type != RBRACKET {
		value := p.parseValue()
		arr = append(arr, value)

		if p.curToken.Type == COMMA {
			p.nextToken()
		}
	}

	p.nextToken() // skip closing bracket
	return arr
}
