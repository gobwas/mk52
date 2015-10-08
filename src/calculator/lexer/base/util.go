package base

func isSymbolic(char rune) bool {
	return char == '$' || char == '_'
}

func isLowerCase(char rune) bool {
	return char >= 'a' && char <= 'z'
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\n' || char == '\r' || char == '\t'
}

func isNumeric(char rune) bool {
	return char >= '0' && char <= '9'
}