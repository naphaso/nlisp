package lexer

import "fmt"

type LexError struct {
	Line   int64
	Column int64
	Desc   string
	Parent error
}

func NewLexError(line, column int64, desc string) error {
	return &LexError{
		Line:   line,
		Column: column,
		Desc:   desc,
	}
}

func NewLexErrorf(line, column int64, format string, a ...interface{}) error {
	return &LexError{
		Line:   line,
		Column: column,
		Desc:   fmt.Sprintf(format, a...),
	}
}

func WrapLexError(line, column int64, desc string, err error) error {
	return &LexError{
		Line:   line,
		Column: column,
		Desc:   desc,
		Parent: err,
	}
}

func (s *LexError) Error() string {
	if s.Parent == nil {
		return fmt.Sprintf("lex error: line %d, col %d: %s", s.Line, s.Column, s.Desc)
	} else {

		return fmt.Sprintf("lex error: line %d, col %d: %s: %s", s.Line, s.Column, s.Desc, s.Parent.Error())
	}
}
