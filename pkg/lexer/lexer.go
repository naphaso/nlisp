package lexer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
)

type lexerState int

const (
	lexerStateNormal lexerState = iota
	lexerStateComment
	lexerStateString
	lexerStateStringEscape
	lexerStateQuote
	lexerStateBacktick
	lexerStateComma
)

var (
	// TODO: extend literal definitions
	regexInteger = regexp.MustCompile("^-?[0-9]+$")
	regexSymbol  = regexp.MustCompile("^[a-z][a-z0-9]*$")
)

func New(r io.RuneScanner) *lexer {
	return &lexer{
		input: r,
		line:  1,
	}
}

type lexer struct {
	input     io.RuneScanner
	nextInput []io.RuneScanner
	tokens    []Token
	buffer    bytes.Buffer

	// state
	endReached bool
	state      lexerState
	line       int64
	column     int64
}

func (s *lexer) PeekToken() (Token, error) {
	for len(s.tokens) == 0 {
		r, _, err := s.input.ReadRune()
		if err != nil {
			if err == io.EOF {
				if len(s.nextInput) > 0 {
					s.input = s.nextInput[0]
					s.nextInput = s.nextInput[1:]
					return s.PeekToken()
				}

				if !s.endReached {
					r = 0
					s.endReached = true
				} else {
					return Token{}, err
				}
			} else {
				return Token{}, err
			}
		}

		err = s.parseRune(r)
		if err != nil {
			return Token{}, err
		}
	}

	return s.tokens[0], nil
}

func (s *lexer) PushData(scanner io.RuneScanner) {
	s.nextInput = append(s.nextInput, scanner)
}

func (s *lexer) GetToken() (Token, error) {
	token, err := s.PeekToken()
	if err != nil {
		return Token{}, err
	}

	s.tokens = s.tokens[1:]
	return token, nil
}

func (s *lexer) parseRune(r rune) error {
	s.column += 1
	var err error
	switch s.state {
	case lexerStateNormal:
		switch r {
		case '(':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.tokens = append(s.tokens, Token{Type: TokenLPar})
		case ')':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.tokens = append(s.tokens, Token{Type: TokenRPar})
		case '\'':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.state = lexerStateQuote
		case '`':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.state = lexerStateBacktick
		case ',':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.state = lexerStateComma
		case '\n':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			fallthrough
		case 0:
			err = s.parseBuffer()
			if err != nil {
				return err
			}
		case ';':
			err = s.parseBuffer()
			if err != nil {
				return err
			}
			s.state = lexerStateComment
		case '"':
			if s.buffer.Len() > 0 {
				return NewLexError(s.line, s.column, "unexpected quote")
			}

			s.state = lexerStateString
		default:
			_, _ = s.buffer.WriteRune(r)
		}
	case lexerStateQuote:
		switch r {
		case '(':
			if s.buffer.Len() > 0 {
				err = s.parseQuotedBuffer()
				if err != nil {
					return err
				}

				s.tokens = append(s.tokens, Token{Type: TokenLPar})
			} else {
				s.tokens = append(s.tokens, Token{Type: TokenQLPar})
			}
			s.state = lexerStateNormal
		case ')':
			if s.buffer.Len() > 0 {
				err = s.parseQuotedBuffer()
				if err != nil {
					return err
				}
				s.state = lexerStateNormal
				s.tokens = append(s.tokens, Token{Type: TokenRPar})
			} else {
				return errors.New("quoted rpar")
			}
		case '\n':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			fallthrough
		case 0:
			if s.buffer.Len() > 0 {
				err = s.parseQuotedBuffer()
				if err != nil {
					return err
				}
				s.state = lexerStateNormal
			} else {
				return errors.New("invalid quotation")
			}
		case ';':
			if s.buffer.Len() > 0 {
				err = s.parseQuotedBuffer()
				if err != nil {
					return err
				}
				s.state = lexerStateNormal
			} else {
				return errors.New("invalid quotation")
			}
		default:
			_, _ = s.buffer.WriteRune(r)
		}
	case lexerStateBacktick:
		switch r {
		case '(':
			if s.buffer.Len() > 0 {
				return errors.New("backtick error")
			}

			s.tokens = append(s.tokens, Token{Type: TokenBLPar})
			s.state = lexerStateNormal
		default:
			return errors.New("unexpected backtick")
		}
	case lexerStateComma:
		switch r {
		case '(':
			if s.buffer.Len() > 0 {
				err = s.parseCommaBuffer()
				if err != nil {
					return err
				}
				s.tokens = append(s.tokens, Token{Type: TokenLPar})
			} else {
				s.tokens = append(s.tokens, Token{Type: TokenCLPar})
			}
			s.state = lexerStateNormal
		case ')':
			if s.buffer.Len() > 0 {
				err = s.parseCommaBuffer()
				if err != nil {
					return err
				}

				s.state = lexerStateNormal
				s.tokens = append(s.tokens, Token{Type: TokenRPar})
			} else {
				return errors.New("comma rpar")
			}
		case '\n':
			fallthrough
		case '\r':
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			fallthrough
		case 0:
			if s.buffer.Len() > 0 {
				err = s.parseCommaBuffer()
				if err != nil {
					return err
				}
				s.state = lexerStateNormal
			} else {
				return errors.New("invalid comma")
			}
		case ';':
			if s.buffer.Len() > 0 {
				err = s.parseCommaBuffer()
				if err != nil {
					return err
				}
				s.state = lexerStateNormal
			} else {
				return errors.New("invalid comma")
			}
		default:
			_, _ = s.buffer.WriteRune(r)
		}
	case lexerStateComment:
		if r == '\n' {
			s.state = lexerStateNormal
			// skip comments
		}
	case lexerStateString:
		if r == '"' {
			s.tokens = append(s.tokens, Token{Type: TokenString, Value: s.buffer.String()})
			s.buffer.Reset()
			s.state = lexerStateNormal
		} else if r == '\\' {
			s.state = lexerStateStringEscape
		} else {
			s.buffer.WriteRune(r)
		}
	case lexerStateStringEscape:
		r, err = s.unescape(r)
		if err != nil {
			return err
		}
		s.buffer.WriteRune(r)
		s.state = lexerStateString
	}

	if r == '\n' {
		s.line += 1
		s.column = 0
	}

	return nil
}

func (s *lexer) parseBuffer() error {
	if s.buffer.Len() > 0 {
		value := s.buffer.String()
		s.buffer.Reset()
		if value == "true" || value == "false" {
			s.tokens = append(s.tokens, Token{
				Type:  TokenBool,
				Value: value,
			})
		} else if value == "." {
			s.tokens = append(s.tokens, Token{
				Type: TokenDot,
			})
		} else if value == "nil" {
			s.tokens = append(s.tokens, Token{
				Type: TokenNil,
			})
		} else if regexInteger.MatchString(value) {
			s.tokens = append(s.tokens, Token{
				Type:  TokenInteger,
				Value: value,
			})
		} else if regexSymbol.MatchString(value) {
			s.tokens = append(s.tokens, Token{
				Type:  TokenSymbol,
				Value: value,
			})
		} else {
			return errors.New(fmt.Sprintf("unknown token: %s", value))
		}
	}

	return nil
}

func (s *lexer) parseQuotedBuffer() error {
	if s.buffer.Len() > 0 {
		value := s.buffer.String()
		s.buffer.Reset()
		if regexSymbol.MatchString(value) {
			s.tokens = append(s.tokens, Token{
				Type:  TokenQSymbol,
				Value: value,
			})
		} else {
			return errors.New("unknown quoted token: " + value)
		}
	}

	return nil
}

func (s *lexer) parseCommaBuffer() error {
	if s.buffer.Len() > 0 {
		value := s.buffer.String()
		s.buffer.Reset()
		if regexSymbol.MatchString(value) {
			s.tokens = append(s.tokens, Token{
				Type:  TokenCSymbol,
				Value: value,
			})
		} else {
			return errors.New("unknown comma token: " + value)
		}
	}

	return nil
}

func (s *lexer) makeAtomErr() error {
	return NewLexErrorf(s.line, s.column-int64(s.buffer.Len()), "failed to parse atom: %s", s.buffer.String())
}
