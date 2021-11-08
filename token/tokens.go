package token

// Type is the smallest unit of the grammer
type Type int

// Token types
const (
	LeftBracket Type = iota
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

	Print
	Get
	Post
	Return

	Nil
	True
	False
	String
	Number

	Newline
	EOF
)
