package sexp

import (
	"errors"
)

type Sexp interface {
	SexpString() string
	Equal(o Sexp) bool
}

func IsList() bool {
	panic("not implemented")
}

func ToArray(e Sexp) ([]Sexp, error) {
	var el []Sexp
	var curr Pair
	var ok bool
	for e != Nil {
		curr, ok = e.(Pair)
		if !ok {
			return nil, errors.New("not a list")
		}

		el = append(el, curr.Head)
		e = curr.Tail
	}

	return el, nil
}

func NewList(elements []Sexp) Sexp {
	root := Nil
	for i := len(elements) - 1; i >= 0; i-- {
		root = NewPair(elements[i], root)
	}
	return root
}

func List(args ...Sexp) Sexp {
	return NewList(args)
}

func NewListBuilder() *ListBuilder {
	return &ListBuilder{}
}

type ListBuilder struct {
	list []Sexp
}

func (s *ListBuilder) Append(v Sexp) {
	s.list = append(s.list, v)
}

func (s *ListBuilder) Build() Sexp {
	// probably there is better way to build this list?
	var root Sexp = Nil
	for i := len(s.list) - 1; i >= 0; i-- {
		root = NewPair(s.list[i], root)
	}
	return root
}

func ExtractTwoArgs(args Sexp) (Sexp, Sexp, error) {
	p, ok := args.(Pair)
	if !ok {
		return nil, nil, errors.New("invalid arguments")
	}

	a1 := p.Head

	p, ok = p.Tail.(Pair)
	if !ok {
		return nil, nil, errors.New("invalid arguments")
	}

	a2 := p.Head

	return a1, a2, nil
}

func ExtractThreeArgs(args Sexp) (Sexp, Sexp, Sexp, error) {
	p, ok := args.(Pair)
	if !ok {
		return nil, nil, nil, errors.New("invalid arguments")
	}

	a1 := p.Head

	p, ok = p.Tail.(Pair)
	if !ok {
		return nil, nil, nil, errors.New("invalid arguments")
	}

	a2 := p.Head

	p, ok = p.Tail.(Pair)
	if !ok {
		return nil, nil, nil, errors.New("invalid arguments")
	}

	a3 := p.Head

	return a1, a2, a3, nil
}
