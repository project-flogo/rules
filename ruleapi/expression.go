package ruleapi

// Type is a type of byte code
type Type uint8

const (
	// TypeLiteral a literal type
	TypeLiteral Type = iota
	// TypeNot a logical not operation
	TypeNot
	// TypeOr a logical or operation
	TypeOr
	// TypeAnd a logical and operation
	TypeAnd
	// TypeEq == operator
	TypeEq
	// TypeNEq != operator
	TypeNEq
	// TypeGt > operator
	TypeGt
	// TypeGte >= operator
	TypeGte
	// TypeLt < operator
	TypeLt
	// TypeLte <= operator
	TypeLte
)

// ByteCode is a expression byte code
type ByteCode struct {
	T       Type
	Literal string
}

// String prints out the string representation of the code
func (code *ByteCode) String() string {
	switch code.T {
	case TypeLiteral:
		return code.Literal
	case TypeNot:
		return "!"
	case TypeOr:
		return "||"
	case TypeAnd:
		return "&&"
	case TypeEq:
		return "=="
	case TypeNEq:
		return "!="
	case TypeGt:
		return ">"
	case TypeGte:
		return ">="
	case TypeLt:
		return "<"
	case TypeLte:
		return "<="
	}
	return ""
}

// Expression is a decision table expression
type Expression struct {
	Code []ByteCode
	Top  int
}

// Init initializes the expression
func (e *Expression) Init(expression string) {
	e.Code = make([]ByteCode, len(expression))
}

// AddOperator adds an operator to the stack
func (e *Expression) AddOperator(operator Type) {
	code, top := e.Code, e.Top
	e.Top++
	code[top].T = operator
}

// AddLiteral adds a literal to the stack
func (e *Expression) AddLiteral(literal string) {
	code, top := e.Code, e.Top
	e.Top++
	code[top].T = TypeLiteral
	code[top].Literal = literal
}

// Evaluate evaluates the expression
func (e *Expression) Evaluate(lhsToken string) string {
	stack, top := make([]string, len(e.Code)), 0
	for _, code := range e.Code[0:e.Top] {
		switch code.T {
		case TypeLiteral:
			stack[top] = code.Literal
			top++
		case TypeNot:
			a := &stack[top-1]
			*a = "!(" + *a + ")"
		case TypeOr:
			a, b := &stack[top-2], &stack[top-1]
			top--
			*a = "(" + *a + " || " + *b + ")"
		case TypeAnd:
			a, b := &stack[top-2], &stack[top-1]
			top--
			*a = "(" + *a + " && " + *b + ")"
		case TypeEq:
			a := &stack[top-1]
			*a = lhsToken + " == " + *a
		case TypeNEq:
			a := &stack[top-1]
			*a = lhsToken + " != " + *a
		case TypeGt:
			a := &stack[top-1]
			*a = lhsToken + " > " + *a
		case TypeGte:
			a := &stack[top-1]
			*a = lhsToken + " >= " + *a
		case TypeLt:
			a := &stack[top-1]
			*a = lhsToken + " < " + *a
		case TypeLte:
			a := &stack[top-1]
			*a = lhsToken + " <= " + *a
		}
	}
	return stack[0]
}
