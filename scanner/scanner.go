package scanner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kingzbauer/scraperlang/token"
)

// Token is a lexer/scanner output
type Token struct {
	Type    token.Type
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
}

// Tokens is a slice of tokens
type Tokens []*Token

func (t Tokens) String() string {
	buf := &strings.Builder{}
	for _, t := range t {
		buf.WriteString(fmt.Sprintf("%v", t))
		buf.WriteByte(' ')
		if t.Type == token.Newline {
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

func (t Token) String() string {
	return fmt.Sprintf("Token[%d:%d]{Type: %s, Lexeme: %v, Literal: %v}", t.Line, t.Column, t.Type, t.Lexeme, t.Literal)
}

// ScannerError returned when the scanner encounters an unexpected character
type ScannerError struct {
	Line, Column int
	Msg          string
}

func (err ScannerError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", err.Line, err.Column, err.Msg)
}

var keywords = map[string]token.Type{
	"true":  token.True,
	"false": token.False,
	"nil":   token.Nil,
}

// Scanner given a byte string will go through each byte character and tokenize them
type Scanner struct {
	start   int
	current int
	line    int
	column  int
	src     []byte
	length  int
	tokens  Tokens
}

// NewScanner initializes a new scanner
func NewScanner(src []byte) *Scanner {
	return &Scanner{src: src, length: len(src)}
}

// ScanTokens goes through the provided src string and performs lexing
func (s *Scanner) ScanTokens() (tokens Tokens, err error) {
	defer func() {
		if val := recover(); val != nil {
			err = val.(error)
		}
	}()

	for {
		if s.isAtEnd() {
			break
		}
		s.start = s.current
		s.scanToken()
	}

	eof := &Token{
		Type:   token.EOF,
		Line:   s.line,
		Column: s.column,
	}
	s.tokens = append(s.tokens, eof)

	tokens = s.tokens
	return tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.length
}

func (s *Scanner) advance() byte {
	char := s.src[s.current]
	s.current++
	return char
}

func (s *Scanner) previous() byte {
	return s.src[s.current-1]
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '[':
		s.add(token.LeftBracket, "[")
		s.column++
	case ']':
		s.add(token.RightBracket, "]")
		s.column++
	case '(':
		s.add(token.LeftParen, "(")
		s.column++
	case ')':
		s.add(token.RightParen, ")")
		s.column++
	case '{':
		s.add(token.LeftCurlyBracket, "{")
		s.column++
	case '}':
		s.add(token.RightCurlyBracket, "}")
		s.column++
	case ',':
		s.add(token.Comma, ",")
		s.column++
	case '.':
		s.add(token.Period, ".")
		s.column++
	case ':':
		s.add(token.Colon, ":")
		s.column++
	case '~':
		s.add(token.Tilde, "~")
		s.column++
	case '=':
		s.add(token.Equal, "=")
		s.column++
	case '\'':
		s.addString('\'')
	case '"':
		s.addString('"')
	case '@':
		s.identifier()
	case '-':
		if s.peek() == '>' {
			s.advance()
			s.add(token.Arrow, "->")
			s.column += 2
		} else {
			s.add(token.Minus, "-")
		}
	case ' ':
		s.column++
	case '\n':
		s.add(token.Newline, "\n")
		s.column = 0
		s.line++
	default:
		if s.isAlpha(char) {
			s.identifier()
		} else if s.isDigit(char) {
			s.number()
		} else {
			panic(ScannerError{
				Line:   s.line,
				Column: s.column,
				Msg:    fmt.Sprintf("encountered unexpected token %q", char),
			})
		}
	}
}

func (s *Scanner) add(typ token.Type, lexeme string, literal ...interface{}) {
	t := &Token{
		Type:   typ,
		Column: s.column,
		Line:   s.line,
		Lexeme: lexeme,
	}
	if len(literal) > 0 {
		t.Literal = literal[0]
	}
	s.addToken(t)
}

func (s *Scanner) addToken(t *Token) {
	s.tokens = append(s.tokens, t)
}

func (s *Scanner) addString(delimiter byte) {
	for !s.isAtEnd() {
		char := s.advance()
		if char == '\n' {
			panic(ScannerError{
				Line:   s.line,
				Column: s.column,
				Msg:    "multiline strings not supported",
			})
		}
		if char == delimiter && s.src[s.current-2] != '\\' {
			break
		}
	}

	if s.previous() != delimiter {
		panic(ScannerError{
			Line:   s.line,
			Column: s.column,
			Msg:    "unterminated string",
		})
	}
	lexeme := string(s.src[s.start:s.current])
	s.add(token.String, lexeme, lexeme[1:len(lexeme)-1])
	s.column += len(lexeme)
}

func (s *Scanner) isAlpha(char byte) bool {
	return char >= 'A' && char <= 'Z' ||
		char >= 'a' && char <= 'z' ||
		char == '_'
}

func (s *Scanner) isAlphaNumeric(char byte) bool {
	return s.isAlpha(char) || s.isDigit(char)
}

func (s *Scanner) isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) identifier() {
	for !s.isAtEnd() && s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	lexeme := string(s.src[s.start:s.current])
	if keyword, ok := keywords[lexeme]; ok {
		switch keyword {
		case token.Nil:
			s.add(token.Nil, lexeme, nil)
		case token.True:
			s.add(token.True, lexeme, true)
		case token.False:
			s.add(token.False, lexeme, false)
		}
	} else {
		// Check if it's a tag
		if s.src[s.start] == '@' {
			s.add(token.Tag, lexeme, lexeme[1:])
		} else {
			s.add(token.Ident, lexeme)
		}
	}

	s.column += len(lexeme)
}

func (s *Scanner) number() {
	for !s.isAtEnd() && s.isDigit(s.peek()) {
		s.advance()
	}

	if !s.isAtEnd() && s.peek() == '.' {
		// consume the period
		s.advance()
		// We expect at least one digit charater after a period
		if !s.isDigit(s.peek()) {
			panic(ScannerError{
				Line:   s.line,
				Column: s.column + s.current - s.start,
				Msg:    "expects a fraction value after period",
			})
		}
		for !s.isAtEnd() && s.isDigit(s.peek()) {
			s.advance()
		}
	}

	lexeme := string(s.src[s.start:s.current])
	literal, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		panic(ScannerError{
			Line:   s.line,
			Column: s.column,
			Msg:    err.Error(),
		})
	}
	s.add(token.Number, lexeme, literal)
	s.column += len(lexeme)
}

func (s *Scanner) peek() byte {
	return s.src[s.current]
}
