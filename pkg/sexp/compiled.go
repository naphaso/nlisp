package sexp

import "fmt"

type Compiled struct {
	Code Evaluable
}

type Evaluable interface {
	Eval(env *Env) (Sexp, error)
}

func NewCompiled(code Evaluable) *Compiled {
	return &Compiled{Code: code}
}

func (s *Compiled) Equal(o Sexp) bool {
	return false
}

func (s *Compiled) SexpString() string {
	return fmt.Sprintf("%#v", s.Code)
}
