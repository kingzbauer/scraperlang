package scraperlang

// TokenType is the smallest unit of the grammer
type TokenType int

// Token types
const (
	LeftBracket TokenType = iota
	RightBracket
	LeftParen
	RightParen
	LeftCurlyBracket
	RightCurlyBracket
	Comma
	Period
	Colon
	Tilde
	Equal
	SingleQuote
	DoubleQuote
	Minus
	Arrow

	Ident
	Tag

	Nil
	True
	False
	String
	Number

	Newline
	EOF
)
