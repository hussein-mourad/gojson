package ast

type Stmt interface {
	GetType() string
}

type Document struct {
	Type string
	body Stmt
}

func (d Document) GetType() string {
	return d.Type
}

type Property struct {
	Type  string
	Key   Identifier
	Value Literal
}

func (p Property) GetType() string {
	return p.Type
}

type Object struct {
	Type    string
	members []Stmt
}

func (o Object) GetType() string {
	return o.Type
}

type Array struct {
	Type     string
	elements []Stmt
}

func (a Array) GetType() string {
	return a.Type
}

// Literals

type Literal interface {
	Stmt
	GetValue() any
}

type StringLiteral struct {
	Type  string
	Value string
}

func (s StringLiteral) GetType() string {
	return s.Type
}

func (s StringLiteral) GetValue() string {
	return s.Value
}

type NumberLiteral struct {
	Type  string
	Value float64
}

func (n NumberLiteral) GetType() string {
	return n.Type
}

func (n NumberLiteral) GetValue() float64 {
	return n.Value
}

type BooleanLiteral struct {
	Type  string
	Value bool
}

func (b BooleanLiteral) GetType() string {
	return b.Type
}

func (b BooleanLiteral) GetValue() bool {
	return b.Value
}

type Identifier struct {
	Type  string
	Value string
}

func (i Identifier) GetType() string {
	return i.Type
}

func (i Identifier) GetValue() string {
	return i.Value
}
