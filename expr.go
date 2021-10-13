package scraperlang

// Environment provides the API methods to access the variables for different scopes
type Environment interface {
	Get(Token) interface{}
	Set(Token, interface{})
}

// Visitor implements the visitor pattern interface
type Visitor interface {
	VisitTaggedClosure(TaggedClosure, Environment) interface{}
	VisitGetExpr(GetExpr, Environment) interface{}
	VisitPrintExpr(PrintExpr, Environment) interface{}
	VisitAssignExpr(AssignExpr, Environment) interface{}
	VisitCallExpr(CallExpr, Environment) interface{}
	VisitClosureExpr(ClosureExpr, Environment) interface{}
	VisitAccessExpr(AccessExpr, Environment) interface{}
	VisitHTMLAttrAccessor(HTMLAttrAccessor, Environment) interface{}
	VisitArrayExpr(ArrayExpr, Environment) interface{}
	VisitMapExpr(MapExpr, Environment) interface{}
	VisitLiteralExpr(LiteralExpr, Environment) interface{}
	VisitIdentExpr(IdentExpr, Environment) interface{}
}

// Expr every expression type must implement the expression interface
type Expr interface {
	Accept(Visitor, Environment) interface{}
}

// TaggedClosure defines a top level closure which can be identifiable by a name
type TaggedClosure struct {
	Name  *Token
	Exprs []Expr
}

// GetExpr use to invoke the http get for the provided url(s)
type GetExpr struct {
	Tag    *Token
	URL    Expr
	Header *Expr
}

// PrintExpr prints the provided arguments
type PrintExpr struct {
	Args []Expr
}

// AssignExpr assigns an expression result to a variable
type AssignExpr struct {
	Name  *Token
	Value Expr
}

// CallExpr invokes a callable with the provided arguments
type CallExpr struct {
	Name      *Token
	Arguments []Expr
}

// ClosureExpr is an untagged closure. In contract to a tagged closure which can only appear
// at the top level score.
// This specific closure cannot appear on the top level definition
type ClosureExpr struct {
	Params Tokens
	Exprs  []Expr
}

// AccessExpr allows to access fields of any object that implements the Getter interface
type AccessExpr struct {
	Var   *Token
	Field *Token
}

// HTMLAttrAccessor allows to retrieve attributes of a Node object
type HTMLAttrAccessor struct {
	Var  *Token
	Attr *Token
}

// ArrayExpr initializes an array
type ArrayExpr struct {
	Entries []Expr
}

// MapExpr initializes a map
type MapExpr struct {
	Keys   Tokens
	Values []Expr
}

// LiteralExpr represents a literal value
type LiteralExpr struct {
	Value *Token
}

// IdentExpr defines a variable in a scope
type IdentExpr struct {
	Name *Token
}
