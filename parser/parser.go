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
		prefix = fmt.Sprintf("[%d:%d]", err.token.Line+1, err.token.Column)
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

// HasErrs checks to see whether we have any parser errors
func (p *Parser) HasErrs() bool {
	return len(p.errs) > 0
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
	exprs := []Expr{}
	for !p.match(token.EOF) {
		closure := p.taggledClosure()
		exprs = append(exprs, closure)
		p.eatAll(token.Newline)
	}

	return exprs
}

func (p *Parser) taggledClosure() Expr {
	taggedClosure := TaggedClosure{}

	p.eatAll(token.Newline)
	closureName := p.consume("Expected a tagged closure", token.Ident)
	p.consume("Expected '{' to start the closure body", token.LeftCurlyBracket)

	taggedClosure.Name = closureName
	taggedClosure.Body = p.body()

	return taggedClosure
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
			exprs = append(exprs, p.printExpr())
		case token.Ident:
			if p.match(token.Equal) {
				// Process an assignment
				exprs = append(exprs, p.assignExpr(t))
			} else if p.match(token.LeftParen) {
				// Process a call expression
				argList := p.expressionList(token.RightParen)
				exprs = append(exprs, CallExpr{Name: LiteralExpr{Value: t}, Arguments: argList})
				p.consume("Call expression requires a closing ')'", token.RightParen)
			} else if p.match(token.Period) {
				// parse an attribute function call
				field := p.consume("Expect an field accessor after '.'", token.Ident)
				var argList []Expr
				if p.match(token.LeftParen) {
					argList = p.expressionList(token.RightParen)
					p.consume("Expect ')' to close functin call", token.RightParen)
				} else {
					// This is a free form function call which needs at least one argument
					argList = p.expressionList()
				}
				accessExpr := AccessExpr{Var: LiteralExpr{Value: t}, Field: field}
				exprs = append(exprs, CallExpr{Name: accessExpr, Arguments: argList})
			} else if p.check(token.Newline) {
				p.addErr(Error{
					msg:   "call expression without parenthesis requires at least one argument",
					token: p.previous(),
				})
				exprs = append(exprs, CallExpr{Name: LiteralExpr{Value: t}})
			} else {
				// This should ideally be a call expression without the parenthesis
				// Requires at least one expression
				argList := p.expressionList()
				if len(argList) == 0 {
					p.addErr(Error{
						msg:   fmt.Sprintf("If this is a function, it requires at least one argument"),
						token: t,
					})
				}
				exprs = append(exprs, CallExpr{Name: LiteralExpr{Value: t}, Arguments: argList})
			}
		default:
			panic(Error{
				token: t,
				msg:   "Invalid expression statement as a top-level statement",
			})
		}
		// Consume a Newline after each expression statement
		p.consume("Expect a 'Newline'", token.Newline)
		p.eatAll(token.Newline)
	}

	p.consume("Expect '}' to close body", token.RightCurlyBracket)
	return exprs
}

func (p *Parser) assignExpr(t *token.Token) Expr {
	value := p.expression()
	return AssignExpr{Name: t, Value: value}
}

func (p *Parser) getExpr(tag ...*token.Token) Expr {
	expr := GetExpr{}
	if len(tag) > 0 {
		expr.Tag = tag[0]
	}

	// We expect at least a single expression as the first argument
	URL := p.expression()
	expr.URL = URL

	// We expect an optional header argument and then a newline to complete the statement
	if !p.check(token.Newline) {
		p.consume("Expect ','", token.Comma)
		httpHeaderExpr := p.expression()
		expr.Header = httpHeaderExpr
	}

	return expr
}

func (p *Parser) printExpr() Expr {
	// We might want to catch any error thrown when parsing the expressions parsed to print statement
	// to give a more meaningful, for now we just allow the normal panic handling at the toplevel parse
	// function
	expr := PrintExpr{}
	// We expect at least one expression
	args := []Expr{p.expression()}
	for p.match(token.Comma) {
		p.eatAll(token.Newline)
		args = append(args, p.expression())
	}
	expr.Args = args

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
	} else if !p.check(
		token.Newline,
		token.Comma,
		token.RightBracket,
		token.RightCurlyBracket,
		token.RightParen) {
		// We need to get an argument list
		argList := p.expressionList()
		expr = CallExpr{Name: expr, Arguments: argList}
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
				}
				key := p.expression()
				p.consume("Expected ']'", token.RightBracket)
				expr = MapAccessExpr{Name: expr, Key: key}
			case token.Period:
				p.advance()
				ident := p.consume("Expect an attribute name after a '.'", token.Ident)
				expr = AccessExpr{Var: expr, Field: ident}
			default:
				return expr
			}
		}
	}
}

// expressionList returns 0 or more expressions separated with a comma
// We use the delimiter token to know if we need to return an empty list if
// the delimiter is encountered as the first thing
//
// TODO: we might want to handle consuming the delimiter here if it's provided
func (p *Parser) expressionList(delimiter ...token.Type) []Expr {
	// Empty expression list
	if p.check(delimiter...) {
		return nil
	}

	// We check if we have a new line.
	// This means that we have an empty expression list without a delimiter which should be an error
	if p.check(token.Newline) {
		p.addErr(Error{
			msg:   "A free form call expression requires at least a single argument",
			token: p.peek(),
		})
		return nil
	}

	exprs := []Expr{p.expression()}

	for p.match(token.Comma) {
		// We can allow at most  one Newline after a comma
		p.eatAll(token.Newline)
		expr := p.expression()
		exprs = append(exprs, expr)
	}

	return exprs
}

func (p *Parser) mapExpr() Expr {
	// Consume any newlines if any
	p.eatAll(token.Newline)
	entries := make(map[string]Expr)
	// if the next token is not a closing curly bracket, process map entries
	if p.peek().Type != token.RightCurlyBracket {
		key, value := p.mapEntry()
		entries[key.Literal.(string)] = value
		for p.match(token.Comma) {
			p.eatAll(token.Newline)
			key, value = p.mapEntry()
			entries[key.Literal.(string)] = value
		}
		p.eatAll(token.Newline)
	}
	p.consume("expect closing '}'", token.RightCurlyBracket)
	return MapExpr{Entries: entries}
}

func (p *Parser) mapEntry() (*token.Token, Expr) {
	// For now, only keys of type string are allowed
	key := p.consume("expect a key of type 'string'", token.String)
	p.consume("expect ':'", token.Colon)
	value := p.expression()
	return key, value
}

func (p *Parser) arrayExpr() Expr {
	exprs := []Expr{}
	// Consume all possible newlines after the opening square bracket
	p.eatAll(token.Newline)
	// If we don't encounter a closing square bracket, we then expect an expression list
	if !p.match(token.RightBracket, token.EOF) {
		expr := p.expression()
		exprs = append(exprs, expr)
		for p.match(token.Comma) {
			p.eatAll(token.Newline)
			exprs = append(exprs, p.expression())
		}
		p.eatAll(token.Newline)
	}
	// All expressions have been consumed to this point, we therefore expect a closing Right Bracket
	p.consume("expect ']'", token.RightBracket)
	return ArrayExpr{Entries: exprs}
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
	p.consume("A closure requires a body", token.RightParen)
	p.consume("Missing '{' to start the closure body", token.LeftCurlyBracket)
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
