package vm

import (
	"errors"
	"fmt"
	"time"

	"github.com/naphaso/nlisp/pkg/sexp"
	"github.com/naphaso/nlisp/pkg/symbol"
)

func New() *VM {
	return &VM{
		env: sexp.NewEnv(),
	}
}

var (
	symbolPlus = symbol.ByName("plus")
	symbolEval = symbol.ByName("eval")
)

type VM struct {
	env    *sexp.Env
	stack1 Stack
	stack2 Stack
}

func (s *VM) Eval(e sexp.Sexp) (sexp.Sexp, error) {
	var err error
	s.stack1.Push(symbolEval)
	s.stack1.Push(e)

	var c sexp.Sexp
	for {
		fmt.Printf("stack1: %s\n", s.stack1.String())
		fmt.Printf("stack2: %s\n", s.stack2.String())

		for {
			c = s.stack1.Pop()
			if c == nil {
				break
			}

			if c.Equal(symbolPlus) {
				err = s.evalPlus()
				if err != nil {
					return nil, err
				}
			} else if c.Equal(symbolEval) {
				err = s.evalEval()
				if err != nil {
					return nil, err
				}
			}
		}

		s.stack1, s.stack2 = s.stack2, s.stack1
		time.Sleep(5 * time.Second)
	}
}

func (s *VM) evalEval() error {
	e := s.stack1.Pop()
	if e == nil {
		return errors.New("eval nil")
	}

	switch ev := e.(type) {
	case *sexp.Int64:
		s.stack2.Push(e)
	case *sexp.Symbol:
		s.stack2.Push(s.env.Get(ev))
	case *sexp.Pair:
		s.stack2.Push(symbolPlus)
		for e != nil {
			p, ok := e.(*sexp.Pair)
			if !ok {
				return errors.New("not a pair")
			}

			s.stack2.Push(symbolEval)
			s.stack2.Push(p.Head)
			e = p.Tail
		}
	}

	return nil
}

func (s *VM) evalPlus() error {
	var value int64
	for {
		e := s.stack1.Pop()
		if e == nil {
			s.stack2.Push(sexp.NewInt64(value))
			break
		}

		ei, ok := e.(*sexp.Int64)
		if !ok {
			return errors.New("plus not integer")
		}

		value += ei.Value
	}
	return nil
}
