package rpn

import (
	"calculator/lexer"
	"fmt"
	"math"
	"strconv"
)

type Function struct {
	n    int
	fn   func([]float64) (float64, error)
	args []float64
}

func NewFunc(n int, fn func([]float64) (float64, error)) Function {
	return Function{
		n:  n,
		fn: fn,
	}
}

func (self *Function) Add(arg float64) {
	self.args = append(self.args, arg)
}

func (self *Function) Calc() (float64, error) {
	if len(self.args) < self.n {
		return 0, fmt.Errorf("Not complete set of arguments yet")
	}

	return self.fn(self.args)
}

var operators = map[lexer.Type]Function{
	lexer.PLUS: NewFunc(2, func(args []float64) (float64, error) {
		return args[0] + args[1], nil
	}),
	lexer.MINUS: NewFunc(2, func(args []float64) (float64, error) {
		return args[0] - args[1], nil
	}),
	lexer.DIVIDE: NewFunc(2, func(args []float64) (float64, error) {
		return args[0] / args[1], nil
	}),
	lexer.MULTIPLY: NewFunc(2, func(args []float64) (float64, error) {
		return args[0] * args[1], nil
	}),
	lexer.POWER: NewFunc(2, func(args []float64) (float64, error) {
		return math.Pow(args[0], args[1]), nil
	}),
	lexer.NEGATE: NewFunc(1, func(args []float64) (float64, error) {
		return -1 * args[0], nil
	}),
}

type Expression struct {
	rpn       *RPN
	functions map[string]Function
	variables map[string]float64
}

func (self *Expression) AddFunc(name string, fn Function) {
	self.functions[name] = fn
}

func (self *Expression) AddVar(name string, val float64) {
	self.variables[name] = val
}

func (self Expression) Evaluate() (float64, error) {
	var stack []float64

	for _, token := range self.rpn.out.tokens {
		switch {

		case token.Type == lexer.NUMBER:
			// todo this could be optimized at parse stage
			f, err := strconv.ParseFloat(token.Literal, 64)
			if err != nil {
				return 0, fmt.Errorf("could not parse number literal: %s", err)
			}
			stack = append(stack, f)

		case token.Type == lexer.VARIABLE:
			vName := varName(token.Literal)

			if v, ok := self.variables[vName]; !ok {
				return 0, fmt.Errorf("unknown variable %s", vName)
			} else {
				stack = append(stack, v)
			}

		case isOperator(token) || token.Type == lexer.FUNCTION:
			var (
				ok bool
				fn Function
			)

			if fn, ok = self.functions[token.Literal]; !ok {
				fn, ok = operators[token.Type]
			}

			if !ok {
				return 0, fmt.Errorf("unknown operator/function %s", token)
			}

			l := len(stack)

			if fn.n > l {
				return 0, fmt.Errorf("malformed expression: %s", token)
			}

			for _, arg := range stack[l-fn.n:] {
				fn.Add(arg)
			}

			result, err := fn.Calc()
			if err != nil {
				return 0, fmt.Errorf("calculation error: %s", err)
			}

			stack = append(stack[:l-fn.n], result)

		default:
			return 0, fmt.Errorf("cnexpected token %s", token)
		}
	}

	if len(stack) == 1 {
		return stack[0], nil
	} else {
		return 0, fmt.Errorf("malformed expression: %s", stack)
	}
}
