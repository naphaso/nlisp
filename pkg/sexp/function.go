package sexp

import (
	"fmt"
	"reflect"
)

type Func func(Sexp, *Env) (Sexp, error)

func NewFunc(f func(Sexp, *Env) (Sexp, error)) Sexp {
	return Func(f)
}

func (s Func) SexpString() string {
	return fmt.Sprintf("%func(%p)", s)
}

func (s Func) Equal(o Sexp) bool {
	if oc, ok := o.(Func); ok {
		f1 := reflect.ValueOf(oc)
		f2 := reflect.ValueOf(s)
		return f1.Pointer() == f2.Pointer()
	}

	return false
}
