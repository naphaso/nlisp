package sexp

import "fmt"

type Rune struct {
	Value rune
}

func NewRune(r rune) Sexp {
	return &Rune{Value: r}
}

func (s Rune) SexpString() string {
	return fmt.Sprintf("%q", s.Value)
}

func (s Rune) Equal(o Sexp) bool {
	if oc, ok := o.(*Rune); ok {
		return oc.Value == s.Value
	}

	return false
}
