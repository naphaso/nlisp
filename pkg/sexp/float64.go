package sexp

import "strconv"

type Float64 struct {
	Value float64
}

func NewFloat64(v float64) Sexp {
	return &Float64{Value: v}
}

func (s Float64) SexpString() string {
	return strconv.FormatFloat(s.Value, 'f', -1, 64)
}

func (s Float64) Equal(o Sexp) bool {
	if oc, ok := o.(*Float64); ok {
		return oc.Value == s.Value
	}

	return false
}
