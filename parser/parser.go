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
	//p.consume("Expected the closure body", token.Newline, token.RightCurlyBracket)
	//p.eatAll(token.Newline)

	//// If the previous token to be consume is `}` then we are done here
	//// we have an empty closure
	//if p.previous().Type == token.RightCurlyBracket {
	//return TaggedClosure{Name: closureName}
	//}

	//// proceed to consume the body

	return nil
}

func (p *Parser) body() []Expr {
	var exprs []Expr

	// At the moment we expect at least one new line after the opening curly bracket
	p.consume("Expect a Newline after and opening bracket", token.Newline)
	p.eatAll(token.Newline)

	for !p.check(token.RightCurlyBracket, token.EOF) {
		t := p.advance()
		switch t.Type {
		case token.Tag:
			p.consume("Expect a get expression after a tag", token.Get)
			exprs = append(exprs, p.getExpr(t))
		case token.Get:
			exprs = append(exprs, p.getExpr())
		case token.Print:
		case token.Ident:
		default:
			panic(Error{
				token: t,
				msg:   "Invalid expression statement as a top-level statement",
			})
		}
	}

	return exprs
}

func (p *Parser) getExpr(tag ...*token.Token) Expr {
	expr := GetExpr{}
	if len(tag) > 0 {
		expr.Tag = tag[0]
	}

	// We expect at least a single expression as the first argument
	URL := p.expression()
	expr.URL = URL

	return expr
}

func (p *Parser) expression() Expr {
	expr := p.htmlAttrAccessor()
	return expr
}

func (p *Parser) htmlAttrAccessor() Expr {
	expr := p.accessor()
	if p.match(token.Tilde) {
		attr := p.consume("HTML attribute identifier expected", token.Ident)
		return HTMLAttrAccessor{Var: expr, Attr: attr}
	}
	return expr
}

// TODO: Rename this
func (p *Parser) accessor() Expr {
	switch p.peek().Type {
	case token.LeftParen:
		p.advance()
		return p.closure()
	case token.LeftBracket:
		p.advance()
		return p.arrayExpr()
	case token.LeftCurlyBracket:
		p.advance()
		return p.mapExpr()
	default:
		expr := p.primary()
		for {
			switch p.peek().Type {
			case token.LeftParen:
				p.advance()
				arguments := p.expressionList(token.RightParen)
				p.consume("Call expression requires a closing ')'", token.RightParen)
				expr = CallExpr{Name: expr, Arguments: arguments}
			case token.LeftBracket:
				p.advance()
				if p.peek().Type == token.RightBracket {
					panic(Error{
						token: p.advance(),
						msg:   "Missing 'Key' value for accessing the map/array",
					})
					key := p.expression()
					p.consume("Expected ']'", token.RightBracket)
					expr = MapAccessExpr{Name: expr, Key: key}
				}
			case token.Period:
				p.advance()
				ident := p.consume("Expect an attribute name after a '.'", token.Ident)
				expr = AccessExpr{Var: expr, Field: ident}
			default:
				break
			}
		}
		return expr
	}
}

// expressionList returns 0 or more expressions separated with a comma
// We use the delimiter token to know if we need to return an empty list if
// encountered as the first thing
func (p *Parser) expressionList(delimiter token.Type) []Expr {
	// Empty expression list
	if p.match(delimiter) {
		return nil
	}
	exprs := []Expr{p.expression()}

	if p.match(token.Comma) {
		// We can allow at most  one Newline after a comma
		p.eatAll(token.Newline)
		expr := p.expression()
		exprs = append(exprs, expr)
	}

	return exprs
}

func (p *Parser) mapExpr() Expr {
	return nil
}

func (p *Parser) arrayExpr() Expr {
	return nil
}

func (p *Parser) closure() Expr {
	// parameter list
	// If the next token is not a closing paren, we expect a parameter list
	paramList := token.Tokens{}
	if p.peek().Type != token.RightParen {
		t := p.consume("Expect a parameter entry", token.Ident)
		paramList = append(paramList, t)
		for p.match(token.Comma) {
			p.eatAll(token.Newline)
			t := p.consume("Expect a parameter entry", token.Ident)
			paramList = append(paramList, t)
		}
	}
	p.consume("A closure requires a body", token.LeftBracket)
	body := p.body()
	return ClosureExpr{Params: paramList, Body: body}
}

func (p *Parser) primary() Expr {
	t := p.advance()
	switch t.Type {
	case token.Number:
		fallthrough
	case token.String:
		fallthrough
	case token.True:
		fallthrough
	case token.False:
		fallthrough
	case token.Nil:
		fallthrough
	case token.Ident:
		return LiteralExpr{t}
	default:
		panic(Error{
			token: t,
			msg:   "Unexpected token",
		})
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) advance() *token.Token {
	if p.isAtEnd() {
		panic(Error{
			msg: "Unexpected end of file",
		})
	}

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
