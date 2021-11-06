package parser

import (
	"github.com/kingzbauer/scraperlang/token"
)

// Environment provides the API methods to access the variables for different scopes
type Environment interface {
	Get(token.Token) interface{}
	Set(token.Token, interface{})
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
	VisitMapAccessExpr(MapAccessExpr, Environment) interface{}
}

// Expr every expression type must implement the expression interface
type Expr interface {
	Accept(Visitor, Environment) interface{}
}

// TaggedClosure defines a top level closure which can be identifiable by a name
type TaggedClosure struct {
	Name *token.Token
	Body []Expr
}

// Accept implements the Expr interface
func (expr TaggedClosure) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitTaggedClosure(expr, env)
}

// GetExpr use to invoke the http get for the provided url(s)
type GetExpr struct {
	Tag    *token.Token
	URL    Expr
	Header Expr
}

// Accept implements the Expr interface
func (expr GetExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitGetExpr(expr, env)
}

// PrintExpr prints the provided arguments
type PrintExpr struct {
	Args []Expr
}

// Accept implements the Expr interface
func (expr PrintExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitPrintExpr(expr, env)
}

// AssignExpr assigns an expression result to a variable
type AssignExpr struct {
	Name  *token.Token
	Value Expr
}

// Accept implements the Expr interface
func (expr AssignExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitAssignExpr(expr, env)
}

// CallExpr invokes a callable with the provided arguments
type CallExpr struct {
	Name      Expr
	Arguments []Expr
}

// Accept implements the Expr interface
func (expr CallExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitCallExpr(expr, env)
}

// ClosureExpr is an untagged closure. In contract to a tagged closure which can only appear
// at the top level score.
// This specific closure cannot appear on the top level definition
type ClosureExpr struct {
	Params token.Tokens
	Body   []Expr
}

// Accept implements the Expr interface
func (expr ClosureExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitClosureExpr(expr, env)
}

// AccessExpr allows to access fields of any object that implements the Getter interface
type AccessExpr struct {
	Var   Expr
	Field *token.Token
}

// Accept implements the Expr interface
func (expr AccessExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitAccessExpr(expr, env)
}

// MapAccessExpr allows accessing an element from a slice by index of map key
type MapAccessExpr struct {
	Name Expr
	Key  Expr
}

// Accept implements the Expr interface
func (expr MapAccessExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitMapAccessExpr(expr, env)
}

// HTMLAttrAccessor allows to retrieve attributes of a Node object
type HTMLAttrAccessor struct {
	Var  Expr
	Attr *token.Token
}

// Accept implements the Expr interface
func (expr HTMLAttrAccessor) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitHTMLAttrAccessor(expr, env)
}

// ArrayExpr initializes an array
type ArrayExpr struct {
	Entries []Expr
}

// Accept implements the Expr interface
func (expr ArrayExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitArrayExpr(expr, env)
}

// MapExpr initializes a map
type MapExpr struct {
	Entries map[string]Expr
}

// Accept implements the Expr interface
func (expr MapExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitMapExpr(expr, env)
}

// LiteralExpr represents a literal value
type LiteralExpr struct {
	Value *token.Token
}

// Accept implements the Expr interface
func (expr LiteralExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitLiteralExpr(expr, env)
}

// IdentExpr defines a variable in a scope
type IdentExpr struct {
	Name *token.Token
}

// Accept implements the Expr interface
func (expr IdentExpr) Accept(visitor Visitor, env Environment) interface{} {
	return visitor.VisitIdentExpr(expr, env)
}
