package ast

type Stmt interface {
	GetType() string
}

type Document struct {
	Type string
	Body Stmt
}

func (d Document) GetType() string {
	return d.Type
}

func NewDocument() *Document {
	return &Document{Type: "Document"}
}

type Property struct {
	Type  string
	Key   Identifier
	Value Literal
}

func (p Property) GetType() string {
	return p.Type
}

func NewProperty() *Property {
	return &Property{Type: "Property"}
}

type Object struct {
	Type    string
	members []Stmt
}

func (o Object) GetType() string {
	return o.Type
}

func NewObject() *Object {
	return &Object{Type: "Object"}
}

type Array struct {
	Type     string
	elements []Stmt
}

func (a Array) GetType() string {
	return a.Type
}

func NewArray() *Array {
	return &Array{Type: "Array"}
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

func NewStringLiteral() *StringLiteral {
	return &StringLiteral{Type: "String"}
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

func NewNumberLiteral() *NumberLiteral {
	return &NumberLiteral{Type: "Number"}
}

type BooleanLiteral struct {
	Type  string
	Value bool
}

func NewBooleanLiteral() *BooleanLiteral {
	return &BooleanLiteral{Type: "Boolean"}
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

func NewIdentifier() *Identifier {
	return &Identifier{Type: "Identifier"}
}
