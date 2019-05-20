package sexp

import "strconv"

type Int64 int64

func NewInt64(v int64) Sexp {
	return Int64(v)
}

func (s Int64) SexpString() string {
	return strconv.FormatInt(int64(s), 10)
}

func (s Int64) Equal(o Sexp) bool {
	if oc, ok := o.(Int64); ok {
		return oc == s
	}

	return false
}
