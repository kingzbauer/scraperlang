package scraperlang

// Token is the smallest unit of the grammer
type Token int

// Token types
const (
	LeftBracket Token = iota
	RightBracket
	LeftParen
	RightParen
	Comma
	Period
	Colon
	Tilde
	Equal
	SingleQuote
	DoubleQuote

	Ident

	Nil
	True
	False
	String
	Number

	Newline
	EOF
)
