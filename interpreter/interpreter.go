package interpreter

import (
	"errors"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"

	"github.com/kingzbauer/scraperlang/parser"
	"github.com/kingzbauer/scraperlang/token"
)

// ReturnException is used to shortcircuit the flow of the current routine by way of panic
type ReturnException struct {
	Value interface{}
}

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
	wg             *sync.WaitGroup
	pool           *ants.Pool
}

// VisitBodyExpr executes all the expressions in the body expressions
func (i *Interpreter) VisitBodyExpr(expr parser.BodyExpr, e parser.Environment) (val interface{}) {
	defer func() {
		if v := recover(); v != nil {
			if returnExp, ok := v.(ReturnException); ok {
				val = returnExp.Value
			} else {
				panic(v)
			}
		}
	}()

	for _, exp := range expr.Exprs {
		exp.Accept(i, e)
	}
	return
}

// VisitReturnExpr evaluates a return expression
func (i *Interpreter) VisitReturnExpr(expr parser.ReturnExpr, e parser.Environment) interface{} {
	var value interface{}
	if expr.Value != nil {
		value = expr.Value.Accept(i, e)
	}
	panic(ReturnException{
		Value: value,
	})
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

	i.wg = &sync.WaitGroup{}
	var err error
	if i.pool, err = ants.NewPool(10, ants.WithPanicHandler(func(val interface{}) {
		if err, ok := val.(Error); ok {
			fmt.Print(err)
		} else {
			panic(val)
		}
	})); err != nil {
		return nil, err
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

	// Wait for all closures to finish before exiting
	i.wg.Wait()
	return
}

// VisitTaggedClosure visits the tagged closure expression
func (i *Interpreter) VisitTaggedClosure(expr parser.TaggedClosure, e parser.Environment) interface{} {
	expr.Body.Accept(i, e)
	return nil
}

// VisitGetExpr given a get expression, executes the requested http call and calls
// the specified tagged closure
func (i *Interpreter) VisitGetExpr(expr parser.GetExpr, e parser.Environment) interface{} {
	var (
		url string
		ok  bool
	)
	val := expr.URL.Accept(i, e)
	if url, ok = val.(string); !ok {
		panic(Error{
			msg: "'get' expects a URL string as it's 1st argument",
		})
	}

	var headers map[string]interface{}
	if expr.Header != nil {
		val := expr.Header.Accept(i, e)
		mapVal, ok := val.(*Map)
		if !ok {
			panic(Error{
				msg: fmt.Sprintf("'get', requires a map as it's 2nd argument"),
			})
		}
		headers = mapVal.instance
	}
	cfg := getWorkConfig{
		// We will use default as the, well, 'default' tag
		tag:     "default",
		url:     url,
		headers: headers,
	}
	if expr.Tag != nil {
		cfg.tag = expr.Tag.Literal.(string)
	}

	work := i.newGetWork(cfg)
	i.wg.Add(1)
	if err := i.pool.Submit(work); err != nil {
		panic(Error{
			msg: err.Error(),
		})
	}

	return nil
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

// VisitCallExpr executes the callable with the given arguments
func (i *Interpreter) VisitCallExpr(expr parser.CallExpr, e parser.Environment) interface{} {
	val := expr.Name.Accept(i, e)
	callable, ok := val.(Callable)
	if !ok {
		panic(Error{
			msg: fmt.Sprintf("%q is not a callable", val),
		})
	}
	if callable.Arity() != len(expr.Arguments) {
		panic(Error{
			msg: fmt.Sprintf("Expect %d arguments, got %d", callable.Arity(), len(expr.Arguments)),
		})
	}

	args := make([]interface{}, len(expr.Arguments))
	// Eval every argument expression
	for index, exp := range expr.Arguments {
		args[index] = exp.Accept(i, e)
	}

	return callable.Call(args...)
}

// VisitClosureExpr creates a new callable from the closure expression
func (i *Interpreter) VisitClosureExpr(expr parser.ClosureExpr, e parser.Environment) interface{} {
	return NewClosure(e, expr, i)
}

// VisitAccessExpr allows the retrieving of attributes from a runtime instance that implements the Accessor
// interface
func (i *Interpreter) VisitAccessExpr(expr parser.AccessExpr, e parser.Environment) interface{} {
	val := expr.Var.Accept(i, e)
	if accessor, ok := val.(Accessor); ok {
		return accessor.Get(expr.Field.Lexeme)
	}
	panic(Error{
		msg: fmt.Sprintf("%s does not implement the Accessor interface", val),
	})
}

func (i *Interpreter) VisitHTMLAttrAccessor(_ parser.HTMLAttrAccessor, _ parser.Environment) interface{} {
	panic("not implemented") // TODO: Implement
}

// VisitArrayExpr creates a runtime list
func (i *Interpreter) VisitArrayExpr(expr parser.ArrayExpr, e parser.Environment) interface{} {
	a := &Array{entries: make([]interface{}, len(expr.Entries))}
	for index, entry := range expr.Entries {
		a.entries[index] = entry.Accept(i, e)
	}

	return a
}

// VisitMapExpr creates a runtime hash map
func (i *Interpreter) VisitMapExpr(expr parser.MapExpr, e parser.Environment) interface{} {
	m := &Map{instance: make(map[string]interface{})}
	for key, value := range expr.Entries {
		v := value.Accept(i, e)
		m.instance[key] = v
	}
	return m
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

// VisitMapAccessExpr indexes into either a list or hash map
func (i *Interpreter) VisitMapAccessExpr(expr parser.MapAccessExpr, e parser.Environment) interface{} {
	val := expr.Name.Accept(i, e)
	// The value needs to implement the Keyer interface
	if keyer, ok := val.(Keyer); ok {
		return keyer.GetValue(expr.Key.Accept(i, e))
	}
	panic(Error{
		msg: fmt.Sprintf("%s cannot be indexed", val),
	})

}
