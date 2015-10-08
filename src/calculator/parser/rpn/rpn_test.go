package rpn

import (
	"calculator/lexer"
	. "github.com/onsi/gomega"
	"testing"
)

type StubScanner struct {
	tokens []lexer.Token
}

func (self *StubScanner) shift() (token *lexer.Token) {
	if len(self.tokens) == 0 {
		return
	}

	token = &self.tokens[0]
	self.tokens = self.tokens[1:]

	return
}

func (self *StubScanner) Scan() lexer.Token {
	if token := self.shift(); token != nil {
		return *token
	}

	return lexer.Token{lexer.EOF, ""}
}

func TestStack(t *testing.T) {
	RegisterTestingT(t)

	var head *lexer.Token
	stack := Stack{}

	head = stack.shift()
	Expect(head).To(BeNil())

	token := lexer.Token{lexer.EOF, "eof!"}
	stack.push(token)
	Expect(len(stack.tokens)).To(Equal(1))

	head = stack.shift()
	Expect(head).To(Not(BeNil()))
	Expect(head).To(Equal(&token))
}

func TestRPNShouldBeOk(t *testing.T) {
	for _, test := range []struct {
		tokens []lexer.Token
		rpn    string
	}{
		{
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.PLUS, `+`},
				lexer.Token{lexer.NUMBER, `1`},
			},
			`1 1 +`,
		},
		{
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.MULTIPLY, `*`},
				lexer.Token{lexer.NEGATE, `-`},
				lexer.Token{lexer.NUMBER, `1`},
			},
			`1 1 _ *`,
		},
		{
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `3`},
				lexer.Token{lexer.PLUS, `+`},
				lexer.Token{lexer.NUMBER, `4`},
				lexer.Token{lexer.MULTIPLY, `*`},
				lexer.Token{lexer.NUMBER, `2`},
			},
			`3 4 2 * +`,
		},
		{
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.PLUS, `+`},
				lexer.Token{lexer.FUNCTION, `max`},
				lexer.Token{lexer.PARENTHESIS_LEFT, `(`},
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.COMMA, `,`},
				lexer.Token{lexer.NUMBER, `2`},
				lexer.Token{lexer.COMMA, `,`},
				lexer.Token{lexer.NUMBER, `3`},
				lexer.Token{lexer.PARENTHESIS_RIGHT, `)`},
			},
			`1 1 2 3 max +`,
		},
		{
			[]lexer.Token{
				lexer.Token{lexer.NUMBER, `3`},
				lexer.Token{lexer.PLUS, `+`},
				lexer.Token{lexer.NUMBER, `4`},
				lexer.Token{lexer.MULTIPLY, `*`},
				lexer.Token{lexer.NUMBER, `2`},
				lexer.Token{lexer.DIVIDE, `/`},
				lexer.Token{lexer.PARENTHESIS_LEFT, `(`},
				lexer.Token{lexer.NUMBER, `1`},
				lexer.Token{lexer.MINUS, `-`},
				lexer.Token{lexer.NUMBER, `5`},
				lexer.Token{lexer.PARENTHESIS_RIGHT, `)`},
				lexer.Token{lexer.POWER, `^`},
				lexer.Token{lexer.NUMBER, `2`},
			},
			`3 4 2 * 1 5 - 2 ^ / +`,
		},
	} {
		parser := New(&StubScanner{test.tokens})
		rpn, err := parser.rpn()
		if err != nil {
			t.Error(err)
			continue
		}

		if rpn.String() != test.rpn {
			t.Errorf("Do not equal:\n\tact: %s\n\texp: %s", rpn, test.rpn)
		}
	}
}
