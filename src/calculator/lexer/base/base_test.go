package base

import (
	"calculator/lexer"
	"reflect"
	"strings"
	"testing"
)

func TestShouldBeOk(t *testing.T) {
	for _, test := range []struct {
		str    string
		tokens []lexer.Token
	}{
		{
			"1",
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `1`},
			},
		},
		{
			"max(1, 2)",
			[]lexer.Token{
				lexer.Token{lexer.FUNCTION, `max`},
				lexer.Token{lexer.PARENTHESIS_LEFT, `(`},
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.COMMA, `,`},
				lexer.Token{lexer.WHITESPACE, ` `},
				lexer.Token{lexer.NUMBER, `2`},
				lexer.Token{lexer.PARENTHESIS_RIGHT, `)`},
			},
		},
		{
			"max(1, $a)",
			[]lexer.Token{
				lexer.Token{lexer.FUNCTION, `max`},
				lexer.Token{lexer.PARENTHESIS_LEFT, `(`},
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.COMMA, `,`},
				lexer.Token{lexer.WHITESPACE, ` `},
				lexer.Token{lexer.VARIABLE, `$a`},
				lexer.Token{lexer.PARENTHESIS_RIGHT, `)`},
			},
		},
		{
			"(1 - 0)",
			[]lexer.Token{
				lexer.Token{lexer.PARENTHESIS_LEFT, `(`},
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.WHITESPACE, ` `},
				lexer.Token{lexer.MINUS, `-`},
				lexer.Token{lexer.WHITESPACE, ` `},
				lexer.Token{lexer.NUMBER, `0`},
				lexer.Token{lexer.PARENTHESIS_RIGHT, `)`},
			},
		},
		{
			"-1",
			[]lexer.Token{
				lexer.Token{lexer.NEGATE, `-`},
				lexer.Token{lexer.NUMBER, `1`},
			},
		},
	} {
		var (
			tokens []lexer.Token
			stop   bool
		)

		scanner := New(strings.NewReader(test.str))

		for !stop {
			if token := scanner.Scan(); token.Type != lexer.EOF {
				tokens = append(tokens, token)
			} else {
				stop = true
			}
		}

		if !reflect.DeepEqual(test.tokens, tokens) {
			t.Errorf("Tokens are not equal:\n\t%s\n\t%s", test.tokens, tokens)
		}
	}
}
