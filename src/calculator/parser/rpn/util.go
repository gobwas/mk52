package rpn

import (
	"calculator/lexer"
	"strings"
)

func isOperator(token lexer.Token) bool {
	switch token.Type {
	case lexer.PLUS, lexer.MINUS, lexer.NEGATE, lexer.DIVIDE, lexer.MULTIPLY, lexer.POWER:
		return true
	default:
		return false
	}
}

func tokenToString(token lexer.Token) string {
	switch token.Type {
	case lexer.NEGATE:
		return `_`
	default:
		return token.Literal
	}
}

func varName(str string) string {
	return strings.TrimPrefix(str, "$")
}
