package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/naphaso/nlisp/pkg/lexer"
	"github.com/naphaso/nlisp/pkg/sexp"
)

type Parser interface {
	Read() (sexp.Sexp, error)
}

type parser struct {
	lex   lexer.Lexer
	stack ListStack
}

func New(lex lexer.Lexer) *parser {
	return &parser{
		lex: lex,
	}
}

func ParseString(v string) (sexp.Sexp, error) {
	p := &parser{
		lex: lexer.New(strings.NewReader(v)),
	}
	return p.Read()
}

func (s *parser) IsEmpty() bool {
	return s.stack.IsEmpty()
}

func (s *parser) Read() (sexp.Sexp, error) {
	for {
		token, err := s.lex.GetToken()
		if err != nil {
			return nil, err
		}

		var ex sexp.Sexp
		var quoted, comma, backquote bool
		switch token.Type {
		case lexer.TokenInteger:
			v, err := strconv.ParseInt(token.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			ex = sexp.NewInt64(v)
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenString:
			ex = sexp.NewString(token.Value)
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenSymbol:
			ex = sexp.NewSymbol(token.Value)
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenQSymbol:
			ex = sexp.List(sexp.NewSymbol("quote"), sexp.NewSymbol(token.Value))
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenCSymbol:
			ex = sexp.List(sexp.NewSymbol("comma"), sexp.NewSymbol(token.Value))
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenBool:
			ex = sexp.NewBool(token.Value == "true")
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenNil:
			ex = nil
			if !s.stack.Add(ex) {
				return ex, nil
			}
		case lexer.TokenDot:
			s.stack.MakePair()
		case lexer.TokenLPar:
			s.stack.Push()
		case lexer.TokenQLPar:
			s.stack.PushQuoted()
		case lexer.TokenCLPar:
			s.stack.PushComma()
		case lexer.TokenBLPar:
			s.stack.PushBackquote()
		case lexer.TokenRPar:
			ex, quoted, comma, backquote, err = s.stack.Pop()
			if err != nil {
				return nil, err
			}

			if quoted {
				ex = sexp.List(sexp.NewSymbol("quote"), ex)
			} else if comma {
				ex = sexp.List(sexp.NewSymbol("comma"), ex)
			} else if backquote {
				ex = sexp.NewPair(sexp.NewSymbol("backquote"), ex)
			}

			if !s.stack.Add(ex) {
				return ex, nil
			}
		default:
			panic(fmt.Sprintf("unknown token: type %d, value: %s", token.Type, token.Value))
		}
	}
}
