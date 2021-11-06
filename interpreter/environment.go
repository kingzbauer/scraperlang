package interpreter

import (
	"fmt"

	"github.com/kingzbauer/scraperlang/parser"
	"github.com/kingzbauer/scraperlang/token"
)

type environment struct {
	entries map[string]interface{}
	parent  parser.Environment
}

// NewEnvironment returns a new Environment implementation with an initial set of values
func NewEnvironment(init map[string]interface{}, parent parser.Environment) parser.Environment {
	e := &environment{entries: make(map[string]interface{}), parent: parent}
	for key, value := range init {
		e.entries[key] = value
	}
	return e

}

// Get checks and returns the given variable from either itself or the parent.
// If both don't get the variable, it panics with undefined variable
func (e *environment) Get(ident token.Token) interface{} {
	if val, found := e.entries[ident.Lexeme]; found {
		return val
	} else if e.parent != nil {
		return e.parent.Get(ident)
	}
	panic(Error{msg: fmt.Sprintf("Undefined variable %q", ident.Lexeme), token: &ident})
}

func (e *environment) Set(ident token.Token, val interface{}) {
	e.entries[ident.Lexeme] = val
}
