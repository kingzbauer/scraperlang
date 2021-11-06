package token

import (
	"fmt"
	"strconv"
	"strings"
)

// Token is a lexer/scanner output
type Token struct {
	Type    Type
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
		if t.Type == Newline {
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

func (t Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("%v", t.Literal)
	}
	return fmt.Sprintf("%s", t.Lexeme)
}

// Error returned when the scanner encounters an unexpected character
type Error struct {
	Line, Column int
	Msg          string
}

func (err Error) Error() string {
	return fmt.Sprintf("[%d:%d] %s", err.Line, err.Column, err.Msg)
}

var keywords = map[string]Type{
	"true":  True,
	"false": False,
	"nil":   Nil,
	"print": Print,
	"get":   Get,
	"post":  Post,
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
		Type:   EOF,
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
		s.add(LeftBracket, "[")
		s.column++
	case ']':
		s.add(RightBracket, "]")
		s.column++
	case '(':
		s.add(LeftParen, "(")
		s.column++
	case ')':
		s.add(RightParen, ")")
		s.column++
	case '{':
		s.add(LeftCurlyBracket, "{")
		s.column++
	case '}':
		s.add(RightCurlyBracket, "}")
		s.column++
	case ',':
		s.add(Comma, ",")
		s.column++
	case '.':
		s.add(Period, ".")
		s.column++
	case ':':
		s.add(Colon, ":")
		s.column++
	case '~':
		s.add(Tilde, "~")
		s.column++
	case '=':
		s.add(Equal, "=")
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
			s.add(Arrow, "->")
			s.column += 2
		} else {
			s.add(Minus, "-")
		}
	case ' ':
		s.column++
	case '\n':
		s.add(Newline, "\n")
		s.column = 0
		s.line++
	default:
		if s.isAlpha(char) {
			s.identifier()
		} else if s.isDigit(char) {
			s.number()
		} else {
			panic(Error{
				Line:   s.line,
				Column: s.column,
				Msg:    fmt.Sprintf("encountered unexpected token %q", char),
			})
		}
	}
}

func (s *Scanner) add(typ Type, lexeme string, literal ...interface{}) {
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
			panic(Error{
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
		panic(Error{
			Line:   s.line,
			Column: s.column,
			Msg:    "unterminated string",
		})
	}
	lexeme := string(s.src[s.start:s.current])
	s.add(String, lexeme, lexeme[1:len(lexeme)-1])
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
		case Nil:
			s.add(Nil, lexeme, nil)
		case True:
			s.add(True, lexeme, true)
		case False:
			s.add(False, lexeme, false)
		default:
			s.add(keyword, lexeme)
		}
	} else {
		// Check if it's a tag
		if s.src[s.start] == '@' {
			s.add(Tag, lexeme, lexeme[1:])
		} else {
			s.add(Ident, lexeme)
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
			panic(Error{
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
		panic(Error{
			Line:   s.line,
			Column: s.column,
			Msg:    err.Error(),
		})
	}
	s.add(Number, lexeme, literal)
	s.column += len(lexeme)
}

func (s *Scanner) peek() byte {
	return s.src[s.current]
}
