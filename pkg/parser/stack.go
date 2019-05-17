package parser

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

type ListStack struct {
	stack []*ListStackFrame
}

type ListStackFrame struct {
	backquote bool
	quoted    bool
	comma     bool
	pair      bool
	elements  []sexp.Sexp
}

func (s *ListStack) Push() {
	s.stack = append(s.stack, &ListStackFrame{})
}

func (s *ListStack) PushQuoted() {
	s.stack = append(s.stack, &ListStackFrame{quoted: true})
}

func (s *ListStack) PushComma() {
	s.stack = append(s.stack, &ListStackFrame{comma: true})
}

func (s *ListStack) PushBackquote() {
	s.stack = append(s.stack, &ListStackFrame{backquote: true})
}

func (s *ListStack) IsEmpty() bool {
	return len(s.stack) == 0
}

var ErrEmpty = errors.New("list stack is empty")

// TODO: improve signature
func (s *ListStack) Pop() (sexp.Sexp, bool, bool, bool, error) {
	if len(s.stack) == 0 {
		return nil, false, false, false, ErrEmpty
	}

	frame := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return frame.Make(), frame.quoted, frame.comma, frame.backquote, nil
}

func (s *ListStack) Add(v sexp.Sexp) bool {
	if len(s.stack) == 0 {
		return false
	}

	s.stack[len(s.stack)-1].Add(v)
	return true
}

func (s *ListStack) MakePair() {
	if len(s.stack) == 0 {
		return
	}

	s.stack[len(s.stack)-1].MakePair()
}

func (s *ListStackFrame) Add(v sexp.Sexp) {
	s.elements = append(s.elements, v)
}

func (s *ListStackFrame) Make() sexp.Sexp {
	if s.pair && len(s.elements) == 2 {
		return sexp.NewPair(s.elements[0], s.elements[1])
	}

	return sexp.NewList(s.elements)
}

func (s *ListStackFrame) MakePair() error {
	if s.pair == true {
		return errors.New("second dot in the list")
	}

	if len(s.elements) != 1 {
		return errors.New("unknown dot in the list")
	}

	s.pair = true
	return nil
}
