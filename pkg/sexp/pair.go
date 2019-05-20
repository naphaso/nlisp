package sexp

import "strings"

type Pair struct {
	Head Sexp
	Tail Sexp
}

func NewPair(head, tail Sexp) Pair {
	return Pair{Head: head, Tail: tail}
}

func (s Pair) SexpString() string {
	var listElements []string
	var ok bool
	for {
		listElements = append(listElements, s.Head.SexpString())
		if s.Tail == Nil {
			return "(" + strings.Join(listElements, " ") + ")"
		}

		s, ok = s.Tail.(Pair)
		if !ok {
			break
		}
	}

	return "(" + s.Head.SexpString() + " . " + s.Tail.SexpString() + ")"
}

func (s Pair) Equal(o Sexp) bool {
	if oc, ok := o.(Pair); ok {
		return oc.Head.Equal(s.Head) && oc.Tail.Equal(s.Tail)
	}

	return false
}
