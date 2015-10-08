package rpn

import (
	"calculator/lexer"
	"testing"
)

func TestShouldBeOk(t *testing.T) {
	for _, test := range []struct {
		rpn    RPN
		result float64
	}{
		{
			RPN{
				out: Stack{
					tokens: []lexer.Token{
						lexer.Token{lexer.NUMBER, `1`},
						lexer.Token{lexer.NUMBER, `2`},
						lexer.Token{lexer.NUMBER, `3`},
						lexer.Token{lexer.PLUS, `+`},
						lexer.Token{lexer.MINUS, `-`},
					},
				},
			},
			-4,
		},
		{
			RPN{
				out: Stack{
					tokens: []lexer.Token{
						lexer.Token{lexer.NUMBER, `9`},
						lexer.Token{lexer.NEGATE, `_`},
						lexer.Token{lexer.NUMBER, `3`},
						lexer.Token{lexer.NEGATE, `_`},
						lexer.Token{lexer.DIVIDE, `/`},
					},
				},
			},
			3,
		},
		{
			RPN{
				out: Stack{
					tokens: []lexer.Token{
						lexer.Token{lexer.NUMBER, `2`},
						lexer.Token{lexer.NUMBER, `8`},
						lexer.Token{lexer.POWER, `^`},
					},
				},
			},
			256,
		},
	} {
		expr := Expression{rpn: &test.rpn}

		result, err := expr.Evaluate()
		if err != nil {
			t.Fatal(err)
			return
		}

		if result != test.result {
			t.Errorf("Not equal:\n\texp: %f\n\tact: %f", test.result, result)
		}
	}
}
