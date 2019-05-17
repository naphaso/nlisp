package sexp

import "strings"

type Pair struct {
	Head Sexp
	Tail Sexp
}

func NewPair(head, tail Sexp) *Pair {
	return &Pair{Head: head, Tail: tail}
}

func (s *Pair) SexpString() string {
	var listElements []string
	curr := s
	var ok bool
	for {
		if curr == nil {
			return "(" + strings.Join(listElements, " ") + ")"
		}
		listElements = append(listElements, ToString(curr.Head))
		if !IsNil(curr.Tail) {
			curr, ok = curr.Tail.(*Pair)
			if !ok {
				break
			}
		} else {
			curr = nil
		}
	}

	return "(" + ToString(s.Head) + " . " + ToString(s.Tail) + ")"
}

func (s Pair) Equal(o Sexp) bool {
	if oc, ok := o.(*Pair); ok {
		return oc.Head.Equal(s.Head) && oc.Tail.Equal(s.Tail)
	}

	return false
}
