package parser

import (
	"fmt"

	"github.com/kingzbauer/scraperlang/token"
)

// Error represents a parse/syntax error
type Error struct {
	token *token.Token
	msg   string
}

func (err Error) Error() string {
	prefix := ""
	if err.token != nil {
		prefix = fmt.Sprintf("[%d:%d]", err.token.Line, err.token.Column)
	}
	return fmt.Sprintf("%s %s", prefix, err.msg)
}

// Parser builds an AST from the provided tokens
type Parser struct {
	tokens  token.Tokens
	current int
	errs    []error
}

// New creates and returns a new parser
func New(tokens token.Tokens) *Parser {
	return &Parser{tokens: tokens}
}

// Err returns any error if present
func (p *Parser) Err() []error {
	return p.errs
}

func (p *Parser) addErr(err error) {
	p.errs = append(p.errs, err)
}

// Parse processes the tokens and returns an AST which is basically
// a list of expression trees
func (p *Parser) Parse() (ast []Expr, err error) {
	defer func() {
		if val := recover(); val != nil {
			if e, isError := val.(error); isError {
				err = e
			}
		}
	}()
	ast = p.globalDefs()
	return
}

func (p *Parser) globalDefs() []Expr {
	return nil
}

func (p *Parser) taggledClosure() Expr {
	p.eatAll(token.Newline)
	closureName := p.consume("Expected a tagged closure", token.Ident)
	p.consume("Expected '{' to start the closure body", token.LeftCurlyBracket)
	p.consume("Expected the closure body", token.Newline, token.RightCurlyBracket)
	p.eatAll(token.Newline)

	// If the previous token to be consume is `}` then we are done here
	// we have an empty closure
	if p.previous().Type == token.RightCurlyBracket {
		return TaggedClosure{Name: closureName}
	}

	// proceed to consume the body

	return nil
}

func (p *Parser) body() []Expr {
	var exprs []Expr

	for !p.check(token.RightCurlyBracket) {
		current := p.peek()
		if current.Type == token.Tag || current.Type == token.Ident && current.Lexeme == "get" {
			exprs = append(exprs, p.getExpr())
		}
	}

	return exprs
}

func (p *Parser) getExpr() Expr {
	expr := GetExpr{}

	getParams := func() {
		ident := p.consume("Expected a 'get' expression", token.Ident)
		if ident.Lexeme != "get" {
			p.addErr(Error{
				token: ident,
				msg:   "Expected a 'get' expression",
			})
			p.eatUntil(token.Newline)
			if p.isAtEnd() {
				panic(Error{
					token: p.previous(),
					msg:   "Got an unexpected EOF",
				})
			}
			return
		}

		urlExpr := p.expression()
		headerExpr := p.expression()
		expr.URL = urlExpr
		expr.Header = headerExpr
	}

	t := p.consume("Expected a tag or function 'get'", token.Tag, token.Ident)
	if t.Type == token.Tag {
		expr.Tag = t
		getParams()
	} else {
		getParams()
	}

	return expr
}

func (p *Parser) expression() Expr {

	return nil
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) advance() *token.Token {
	t := p.tokens[p.current]
	p.current++
	return t
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

// match returns a boolean indicating whether the current token matches the
// provides token types
func (p *Parser) match(types ...token.Type) bool {
	if p.check(types...) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(types ...token.Type) bool {
	t := p.peek()
	for _, typ := range types {
		if typ == t.Type {
			return true
		}
	}
	return false
}

func (p *Parser) consume(msg string, typ ...token.Type) *token.Token {
	if p.isAtEnd() {
		panic(Error{
			msg: fmt.Sprintf("%s. %s", msg, "Got unexpected EOF"),
		})
	}
	if !p.match(typ...) {
		t := p.peek()
		panic(Error{
			token: t,
			msg:   fmt.Sprintf("%s. Got unexpected %s", msg, t),
		})
	}
	return p.previous()
}

// eatAll consumes all consequetive tokens that match the provided type and stops
// when they find a token of a different type
func (p *Parser) eatAll(typ token.Type) {
	for p.peek().Type == typ {
		p.advance()
	}
}

func (p *Parser) eatUntil(typ token.Type) {
	for !p.check(typ) && !p.isAtEnd() {
		p.advance()
	}
}
