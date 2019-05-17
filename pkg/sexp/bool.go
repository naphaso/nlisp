package sexp

type Bool struct {
	Value bool
}

func NewBool(v bool) Sexp {
	return &Bool{Value: v}
}

func (s Bool) SexpString() string {
	// TODO: symbols t and nil
	if s.Value {
		return "true"
	} else {
		return "false"
	}
}

func (s Bool) Equal(o Sexp) bool {
	if oc, ok := o.(*Bool); ok {
		return oc.Value == s.Value
	}

	return false
}
