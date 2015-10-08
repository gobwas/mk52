package base

import (
	"bufio"
	"bytes"
	"calculator/lexer"
	"io"
)

type Scanner struct {
	reader *bufio.Reader
	prev   lexer.Type
}

const eof = rune(0)

func New(reader io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(reader), lexer.EOF}
}

func (self *Scanner) read() rune {
	r, _, err := self.reader.ReadRune()
	if err != nil {
		return eof
	}

	return r
}

func (self *Scanner) unread() {
	self.reader.UnreadRune()
}

func (self *Scanner) holdPrev(token lexer.Token) lexer.Token {
	switch token.Type {
	case lexer.WHITESPACE, lexer.ILLEGAL, lexer.EOF:
		return token
	default:
		self.prev = token.Type
		return token
	}
}

func (self *Scanner) Scan() (token lexer.Token) {
	char := self.read()

	switch {
	case isLowerCase(char):
		self.unread()
		return self.scanFunction()
	case char == '$':
		self.unread()
		return self.scanVariable()
	case isWhitespace(char):
		self.unread()
		return self.holdPrev(self.scanWhitespace())
	case isNumeric(char):
		self.unread()
		return self.holdPrev(self.scanNumber())
	}

	switch char {
	case '+':
		return lexer.Token{lexer.PLUS, string(char)}
	case '-':
		switch self.prev {
		case lexer.NUMBER, lexer.PARENTHESIS_LEFT, lexer.PARENTHESIS_RIGHT:
			token = lexer.Token{lexer.MINUS, string(char)}
		default:
			token = lexer.Token{lexer.NEGATE, string(char)}
		}
	case '^':
		token = lexer.Token{lexer.POWER, string(char)}
	case '/':
		token = lexer.Token{lexer.DIVIDE, string(char)}
	case '*':
		token = lexer.Token{lexer.MULTIPLY, string(char)}
	case '(':
		token = lexer.Token{lexer.PARENTHESIS_LEFT, string(char)}
	case ')':
		token = lexer.Token{lexer.PARENTHESIS_RIGHT, string(char)}
	case ',':
		token = lexer.Token{lexer.COMMA, string(char)}
	case eof:
		token = lexer.Token{lexer.EOF, ""}
	default:
		token = lexer.Token{lexer.ILLEGAL, string(char)}
	}

	return self.holdPrev(token)
}

func (self *Scanner) scanFunction() lexer.Token {
	return lexer.Token{lexer.FUNCTION, self.readIdent()}
}

func (self *Scanner) scanVariable() lexer.Token {
	return lexer.Token{lexer.VARIABLE, self.readIdent()}
}

func (self *Scanner) readIdent() string {
	var (
		buffer bytes.Buffer
		stop   bool
	)

	for !stop {
		char := self.read()

		if isLowerCase(char) || isNumeric(char) || isSymbolic(char) {
			buffer.WriteRune(char)
			continue
		}

		switch char {
		case eof:
			stop = true
		default:
			stop = true
			self.unread()
		}
	}

	return buffer.String()
}

func (self *Scanner) scanWhitespace() lexer.Token {
	var (
		buffer bytes.Buffer
		stop   bool
	)

	for !stop {
		char := self.read()

		if isWhitespace(char) {
			buffer.WriteRune(char)
			continue
		}

		switch char {
		case eof:
			stop = true
		default:
			stop = true
			self.unread()
		}
	}

	return lexer.Token{lexer.WHITESPACE, buffer.String()}
}

func (self *Scanner) scanNumber() lexer.Token {
	var (
		buffer bytes.Buffer
		stop   bool
	)

	for !stop {
		char := self.read()

		if isNumeric(char) {
			buffer.WriteRune(char)
			continue
		}

		switch char {
		case eof:
			stop = true
		default:
			stop = true
			self.unread()
		}
	}

	return lexer.Token{lexer.NUMBER, buffer.String()}
}
