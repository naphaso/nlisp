package sexp

import "fmt"

type Symbol struct {
	Code uint64
	Name string
}

func NewSymbol(code uint64, name string) *Symbol {
	return &Symbol{
		Code: code,
		Name: name,
	}
}

func (s Symbol) SexpString() string {
	return fmt.Sprintf("%s", s.Name)
}

func (s Symbol) Equal(o Sexp) bool {
	if oc, ok := o.(*Symbol); ok {
		return oc.Code == s.Code
	}

	return false
}
