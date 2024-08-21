package main

import (
	"fmt"
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
}

func (p *Parser) parse() (interface{}, error) {
	if p.curToken.Type == LBRACE {
		return p.parseObject()
	} else {
		return nil, fmt.Errorf("error: expected '{' at line: %v column: %v", p.lexer.line, p.lexer.column)
	}
}

func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case STRING:
		value := p.curToken.Value
		p.nextToken()
		return value, nil
	case NUMBER:
		value := p.curToken.Value
		p.nextToken()
		return value, nil
	case LBRACE:
		obj, err := p.parseObject()
		if err != nil {
			return nil, err
		}
		return obj, nil
	case LBRACKET:
		obj, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		return obj, nil
	case TRUE:
		p.nextToken()
		return true, nil
	case FALSE:
		p.nextToken()
		return false, nil
	case NULL:
		p.nextToken()
		return nil, nil
	}
	return nil, nil
}

func (p *Parser) parseObject() (interface{}, error) {
	obj := make(map[string]interface{})
	p.nextToken() // Skip opening brace
	for p.curToken.Type != RBRACE {
		key := p.curToken.Value
		p.nextToken()

		if p.curToken.Value != ":" {
			// log.Fatalf("error: expected ':' at line: %v column: %v", p.lexer.line, p.lexer.column)
			return nil, fmt.Errorf("error: expected ':' at line: %v column: %v", p.lexer.line, p.lexer.column)
		}
		p.nextToken()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		if p.curToken.Type == COMMA {
			p.nextToken()
		}
	}
	p.nextToken() // Skip closing brace
	return obj, nil
}

func (p *Parser) parseArray() (interface{}, error) {
	var arr []interface{}

	p.nextToken() // skip opening bracket

	for p.curToken.Type != RBRACKET {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		if p.curToken.Type == COMMA {
			p.nextToken()
		}
	}

	p.nextToken() // skip closing bracket
	return arr, nil
}
