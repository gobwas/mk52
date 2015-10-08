package rpn

import (
	"bytes"
	"calculator/lexer"
)

type Stack struct {
	tokens []lexer.Token
}

func (self *Stack) push(token lexer.Token) {
	self.tokens = append(self.tokens, token)
}

func (self *Stack) shift() (token *lexer.Token) {
	if token = self.head(); token != nil {
		self.tokens = self.tokens[1:]
	}

	return
}

func (self *Stack) pop() (token *lexer.Token) {
	if token = self.last(); token != nil {
		self.tokens = self.tokens[0 : len(self.tokens)-1]
	}

	return
}

func (self Stack) head() (token *lexer.Token) {
	if len(self.tokens) == 0 {
		return
	}

	token = &self.tokens[0]

	return
}

func (self Stack) last() (token *lexer.Token) {
	l := len(self.tokens)
	if l == 0 {
		return
	}

	token = &self.tokens[l-1]

	return
}

func (self Stack) String() string {
	var buf bytes.Buffer
	for _, token := range self.tokens {
		buf.WriteString(tokenToString(token))
		buf.WriteRune(' ')
	}

	return string(bytes.TrimSuffix(buf.Bytes(), []byte(` `)))
}
