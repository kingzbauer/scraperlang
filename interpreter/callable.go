package interpreter

import (
	"fmt"

	"github.com/kingzbauer/scraperlang/parser"
)

// Callable is any runtime object that can be executed as a function call
type Callable interface {
	Call(args ...interface{}) interface{}
	Arity() int
}

// closure is the runtime instance of a closure definition
type closure struct {
	closureEnv parser.Environment
	closureAst parser.ClosureExpr
	i          *Interpreter
}

func (c *closure) Call(args ...interface{}) interface{} {
	initEnv := map[string]interface{}{}
	for i, arg := range args {
		initEnv[c.closureAst.Params[i].Lexeme] = arg
	}
	e := NewEnvironment(initEnv, c.closureEnv)
	var val interface{}
	for _, expr := range c.closureAst.Body {
		val = expr.Accept(c.i, e)
	}
	return val
}

func (c *closure) Arity() int {
	return len(c.closureAst.Params)
}

func (c *closure) String() string {
	return fmt.Sprint("#Closure")
}

// NewClosure creates a new callable for the specific AST closure
func NewClosure(closureEnv parser.Environment, ast parser.ClosureExpr, i *Interpreter) Callable {
	return &closure{closureEnv: closureEnv, closureAst: ast, i: i}
}
