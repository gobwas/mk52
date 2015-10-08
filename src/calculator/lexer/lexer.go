package lexer

import "fmt"

type Token struct {
	Type    Type
	Literal string
}

func (self Token) String() string {
	var t string
	switch self.Type {
	case EOF:
		t = "EOF"
	case WHITESPACE:
		t = "WHITESPACE"
	case ILLEGAL:
		t = "ILLEGAL"
	case NUMBER:
		t = "NUMBER"
	case PLUS:
		t = "PLUS"
	case MINUS:
		t = "MINUS"
	case NEGATE:
		t = "NEGATE"
	case MULTIPLY:
		t = "MULTIPLY"
	case DIVIDE:
		t = "DIVIDE"
	case POWER:
		t = "POWER"
	case FUNCTION:
		t = "FUNCTION"
	case VARIABLE:
		t = "VARIABLE"
	case COMMA:
		t = "COMMA"
	case PARENTHESIS_LEFT:
		t = "PARENTHESIS_LEFT"
	case PARENTHESIS_RIGHT:
		t = "PARENTHESIS_RIGHT"
	}

	return fmt.Sprintf("%s[%s]", t, self.Literal)
}

type Type int

const (
	EOF Type = iota
	WHITESPACE
	ILLEGAL
	NUMBER
	PLUS
	MINUS
	NEGATE
	MULTIPLY
	DIVIDE
	POWER
	FUNCTION
	VARIABLE
	COMMA
	PARENTHESIS_LEFT
	PARENTHESIS_RIGHT
)

type Scanner interface {
	Scan() Token
}
