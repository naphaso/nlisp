package sexp

import "fmt"

type String struct {
	Value string
}

func NewString(v string) Sexp {
	return &String{Value: v}
}

func (s String) SexpString() string {
	return fmt.Sprintf("%q", s.Value)
}

func (s String) Equal(o Sexp) bool {
	if oc, ok := o.(*String); ok {
		return oc.Value == s.Value
	}

	return false
}
