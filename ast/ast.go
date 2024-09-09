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
	Key   *Identifier
	Value Stmt
}

func (p Property) GetType() string {
	return p.Type
}

func NewProperty() *Property {
	Identifier := NewIdentifier()
	return &Property{Type: "Property", Key: Identifier}
}

type Object struct {
	Type    string
	Members []*Property
}

func (o Object) GetType() string {
	return o.Type
}

func NewObject() *Object {
	return &Object{Type: "Object", Members: make([]*Property, 0)}
}

type Array struct {
	Type     string
	Elements []Stmt
}

func (a Array) GetType() string {
	return a.Type
}

func NewArray() *Array {
	return &Array{Type: "Array", Elements: make([]Stmt, 0)}
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

func NewStringLiteral(Value string) *StringLiteral {
	return &StringLiteral{Type: "String", Value: Value}
}

type NumberLiteral struct {
	Type  string
	Value interface{}
}

func (n NumberLiteral) GetType() string {
	return n.Type
}

func (n NumberLiteral) GetValue() interface{} {
	return n.Value
}

func NewNumberLiteral(Value interface{}) *NumberLiteral {
	return &NumberLiteral{Type: "Number", Value: Value}
}

type BooleanLiteral struct {
	Type  string
	Value bool
}

func NewBooleanLiteral(Value bool) *BooleanLiteral {
	return &BooleanLiteral{Type: "Boolean", Value: Value}
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
