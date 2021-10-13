// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package scraperlang

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LeftBracket-0]
	_ = x[RightBracket-1]
	_ = x[LeftParen-2]
	_ = x[RightParen-3]
	_ = x[LeftCurlyBracket-4]
	_ = x[RightCurlyBracket-5]
	_ = x[Comma-6]
	_ = x[Period-7]
	_ = x[Colon-8]
	_ = x[Tilde-9]
	_ = x[Equal-10]
	_ = x[SingleQuote-11]
	_ = x[DoubleQuote-12]
	_ = x[Minus-13]
	_ = x[Arrow-14]
	_ = x[Ident-15]
	_ = x[Tag-16]
	_ = x[Nil-17]
	_ = x[True-18]
	_ = x[False-19]
	_ = x[String-20]
	_ = x[Number-21]
	_ = x[Newline-22]
	_ = x[EOF-23]
}

const _TokenType_name = "LeftBracketRightBracketLeftParenRightParenLeftCurlyBracketRightCurlyBracketCommaPeriodColonTildeEqualSingleQuoteDoubleQuoteMinusArrowIdentTagNilTrueFalseStringNumberNewlineEOF"

var _TokenType_index = [...]uint8{0, 11, 23, 32, 42, 58, 75, 80, 86, 91, 96, 101, 112, 123, 128, 133, 138, 141, 144, 148, 153, 159, 165, 172, 175}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
