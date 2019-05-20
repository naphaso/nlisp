package vm

import (
	"errors"
	"strings"

	"github.com/naphaso/nlisp/pkg/sexp"
)

type Stack struct {
	stack []sexp.Sexp
}

func (s *Stack) Push(e sexp.Sexp) {
	s.stack = append(s.stack, e)
}

func (s *Stack) Pop() sexp.Sexp {
	if len(s.stack) == 0 {
		return nil
	}

	head := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return head
}

func (s *Stack) PushList(e sexp.Sexp) error {
	for e != nil {
		p, ok := e.(*sexp.Pair)
		if !ok {
			return errors.New("not a list")
		}

		s.Push(p.Head)
		e = p.Tail
	}
	return nil
}

func (s *Stack) String() string {
	var ss []string
	for _, v := range s.stack {
		ss = append(ss, sexp.ToString(v))
	}

	return "[" + strings.Join(ss, " ") + "]"
}
