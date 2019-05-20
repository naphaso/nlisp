package sexp

import (
	"fmt"
	"reflect"
)

type RawFunc func(Sexp, *Env) (Sexp, error)
type CompileFunc func(Sexp, *Env) (Sexp, error)

type Func struct {
	Raw     RawFunc
	Compile CompileFunc
}

func NewFunc(f RawFunc, c CompileFunc) Sexp {
	return &Func{
		Raw:     f,
		Compile: c,
	}
}

func (s *Func) SexpString() string {
	return fmt.Sprintf("%func(%p)", s)
}

func (s *Func) Equal(o Sexp) bool {
	if oc, ok := o.(*Func); ok {
		f1 := reflect.ValueOf(oc)
		f2 := reflect.ValueOf(s)
		return f1.Pointer() == f2.Pointer()
	}

	return false
}
