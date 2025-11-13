package ast

// Resource represents a cloud resource in the AST
type Resource struct {
	Name       string
	Type       string
	Properties map[string]interface{}
	DependsOn  []string
}

// Program represents the entire tblang program
type Program struct {
	CloudVendors map[string]*CloudVendor
	Variables    map[string]*Variable
	Resources    []*Resource
}

// CloudVendor represents cloud provider configuration
type CloudVendor struct {
	Name       string
	Properties map[string]interface{}
}

// Variable represents a declared variable
type Variable struct {
	Name  string
	Value interface{}
}

// Expression represents any expression in the language
type Expression interface {
	Evaluate() interface{}
}

// StringLiteral represents a string value
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Evaluate() interface{} {
	return s.Value
}

// NumberLiteral represents a numeric value
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) Evaluate() interface{} {
	return n.Value
}

// BooleanLiteral represents a boolean value
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) Evaluate() interface{} {
	return b.Value
}

// ObjectLiteral represents an object/map value
type ObjectLiteral struct {
	Properties map[string]Expression
}

func (o *ObjectLiteral) Evaluate() interface{} {
	result := make(map[string]interface{})
	for key, expr := range o.Properties {
		result[key] = expr.Evaluate()
	}
	return result
}

// ArrayLiteral represents an array value
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) Evaluate() interface{} {
	var result []interface{}
	for _, expr := range a.Elements {
		result = append(result, expr.Evaluate())
	}
	return result
}