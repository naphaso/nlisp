package lexer

func (s *lexer) unescape(r rune) (rune, error) {
	switch r {
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 'a':
		return '\a', nil
	case 't':
		return '\t', nil
	case '\\':
		return '\\', nil
	case '"':
		return '"', nil
	case '\'':
		return '\'', nil
	case '#':
		return '#', nil
	}
	return ' ', NewLexErrorf(s.line, s.column, "invalid escape symbol: %c", r)
}
