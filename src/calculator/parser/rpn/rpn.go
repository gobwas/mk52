package rpn

import (
	"calculator/lexer"
	"fmt"
	"strconv"
)

var precedence = map[lexer.Type]int{
	lexer.PLUS:     1,
	lexer.MINUS:    1,
	lexer.MULTIPLY: 2,
	lexer.DIVIDE:   2,
	lexer.POWER:    3,
	lexer.NEGATE:   4,
}

type Association int

const (
	LEFT Association = iota
	RIGHT
)

var association = map[lexer.Type]Association{
	lexer.PLUS:     LEFT,
	lexer.MINUS:    LEFT,
	lexer.MULTIPLY: LEFT,
	lexer.DIVIDE:   LEFT,
	lexer.POWER:    RIGHT,
	lexer.NEGATE:   RIGHT,
}

type RPN struct {
	stack Stack
	out   Stack
}

func (self RPN) String() string {
	return self.out.String()
}

type Parser struct {
	scanner lexer.Scanner
	buf     struct {
		size  int
		token lexer.Token
	}
}

func New(scanner lexer.Scanner) *Parser {
	return &Parser{scanner: scanner}
}

func (self *Parser) scan() lexer.Token {
	if self.buf.size == 1 {
		self.buf.size = 0
		return self.buf.token
	}

	token := self.scanner.Scan()
	self.buf.token = token

	return token
}

func (self *Parser) unscan() {
	self.buf.size = 1
}

func (self *Parser) scanWithoutWhitespace() lexer.Token {
	token := self.scan()

	if token.Type == lexer.WHITESPACE {
		return self.scan()
	}

	return token
}

func (self *Parser) parseNumber(literal string) (f float64, err error) {
	f, err = strconv.ParseFloat(literal, 64)
	return
}

func (self *Parser) Parse() (*Expression, error) {
	rpn, err := self.rpn()
	if err != nil {
		return nil, err
	}

	return &Expression{rpn: rpn}, nil
}

func (self *Parser) rpn() (*RPN, error) {
	var (
		stop bool
		rpn  RPN
	)

	for !stop {
		token := self.scan()
		switch {

		case token.Type == lexer.NUMBER:
			rpn.out.push(token)

		case token.Type == lexer.VARIABLE:
			rpn.out.push(token)

		case token.Type == lexer.FUNCTION:
			rpn.stack.push(token)

		case token.Type == lexer.COMMA:
			for {
				last := rpn.stack.last()
				if last == nil {
					return nil, fmt.Errorf("Parentheses mismatched")
				}
				if last.Type == lexer.PARENTHESIS_LEFT {
					break
				}

				rpn.out.push(*rpn.stack.pop())
			}

		case token.Type == lexer.PARENTHESIS_LEFT:
			rpn.stack.push(token)

		case token.Type == lexer.PARENTHESIS_RIGHT:
			for {
				token := rpn.stack.pop()
				if token == nil {
					return nil, fmt.Errorf("Parentheses mismatched")
				}
				if token.Type == lexer.PARENTHESIS_LEFT {
					break
				}

				rpn.out.push(*token)
			}

			last := rpn.stack.last()
			if last != nil && last.Type == lexer.FUNCTION {
				rpn.out.push(*rpn.stack.pop())
			}

		case isOperator(token):
			for {
				head := rpn.stack.last()
				if head == nil || !isOperator(*head) {
					break
				}

				switch association[token.Type] {
				case LEFT:
					if precedence[token.Type] <= precedence[head.Type] {
						rpn.out.push(*rpn.stack.pop())
						continue
					}
				case RIGHT:
					if precedence[token.Type] < precedence[head.Type] {
						rpn.out.push(*rpn.stack.pop())
						continue
					}
				}

				break
			}

			rpn.stack.push(token)

		case token.Type == lexer.EOF:
			var finish bool

			for !finish {
				token := rpn.stack.pop()
				if token == nil {
					finish = true
					continue
				}

				rpn.out.push(*token)
			}

			stop = true
		}
	}

	return &rpn, nil
}
