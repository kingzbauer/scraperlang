package interpreter

import (
	"errors"
	"fmt"

	"github.com/kingzbauer/scraperlang/parser"
	"github.com/kingzbauer/scraperlang/token"
)

// Error defines interpreter errors
type Error struct {
	token *token.Token
	msg   string
}

func (err Error) Error() string {
	prefix := ""
	if err.token != nil {
		prefix = fmt.Sprintf("[%d:%d]", err.token.Line+1, err.token.Column)
	}
	return fmt.Sprintf("%s %s", prefix, err.msg)
}

// Interpreter implements the Visitor interface and the Eval loop
type Interpreter struct {
	ast            []parser.Expr
	taggedClosures map[string]parser.TaggedClosure
}

// New creates a new Intepreter instance
func New(ast []parser.Expr) (*Interpreter, error) {
	i := &Interpreter{}
	i.taggedClosures = make(map[string]parser.TaggedClosure)
	// 1. We expect the top level expression to be tagged closures
	// 2. There is 1 required tagged closure: 'init'
	for _, expr := range ast {
		if closure, ok := expr.(parser.TaggedClosure); ok {
			i.taggedClosures[closure.Name.Lexeme] = closure
		} else {
			return nil, errors.New("Only tagged closures are allowed as global variables")
		}
	}

	// Assert that we have the init closure
	if _, ok := i.taggedClosures["init"]; !ok {
		return nil, errors.New("Missing 'init' tagged closure")
	}

	return i, nil
}

// Exec starts the execution flow for the interpreter
func (i *Interpreter) Exec() (err error) {
	defer func() {
		if val := recover(); val != nil {
			if er, ok := val.(error); ok {
				err = er
			} else {
				err = fmt.Errorf("%v", val)
			}
		}
	}()

	e := NewEnvironment(nil, nil)
	// we start our execution from the init closure
	i.taggedClosures["init"].Accept(i, e)

	return
}

// VisitTaggedClosure visits the tagged closure expression
func (i *Interpreter) VisitTaggedClosure(expr parser.TaggedClosure, e parser.Environment) interface{} {
	for _, exp := range expr.Body {
		exp.Accept(i, e)
	}
	return nil
}

func (i *Interpreter) VisitGetExpr(_ parser.GetExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

// VisitPrintExpr prints the provided arguments to stdout
func (i *Interpreter) VisitPrintExpr(expr parser.PrintExpr, e parser.Environment) interface{} {
	values := make([]interface{}, len(expr.Args))
	for index, expr := range expr.Args {
		values[index] = expr.Accept(i, e)
	}

	fmt.Println(values...)
	return nil
}

// VisitAssignExpr creates a new variable with the value as the expression value
func (i *Interpreter) VisitAssignExpr(expr parser.AssignExpr, e parser.Environment) interface{} {
	val := expr.Value.Accept(i, e)
	e.Set(*expr.Name, val)
	return val
}

func (i *Interpreter) VisitCallExpr(_ parser.CallExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

func (i *Interpreter) VisitClosureExpr(_ parser.ClosureExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

func (i *Interpreter) VisitAccessExpr(_ parser.AccessExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

func (i *Interpreter) VisitHTMLAttrAccessor(_ parser.HTMLAttrAccessor, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

func (i *Interpreter) VisitArrayExpr(_ parser.ArrayExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

func (i *Interpreter) VisitMapExpr(_ parser.MapExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

// VisitLiteralExpr returns the underlying literal value
func (i *Interpreter) VisitLiteralExpr(expr parser.LiteralExpr, e parser.Environment) interface{} {
	switch expr.Value.Type {
	case token.String, token.Number, token.Nil, token.True, token.False:
		return expr.Value.Literal
	case token.Ident:
		return e.Get(*expr.Value)
	}
	return nil
}

// VisitIdentExpr accesses a stored value
func (i *Interpreter) VisitIdentExpr(expr parser.IdentExpr, e parser.Environment) interface{} {
	return e.Get(*expr.Name)
}

func (i *Interpreter) VisitMapAccessExpr(_ parser.MapAccessExpr, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}
