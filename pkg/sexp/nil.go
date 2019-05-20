package sexp

type SexpNil struct{}

var Nil Sexp = SexpNil{}

func NewNil() Sexp {
	return Nil
}

func (s SexpNil) SexpString() string {
	return "nil"
}

func (s SexpNil) Equal(o Sexp) bool {
	return o == Nil
}
