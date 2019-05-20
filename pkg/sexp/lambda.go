package sexp

type Lambda struct {
	Args Sexp
	Body Sexp
	Env  *Env
}

func NewLambda(args Sexp, body Sexp, e *Env) *Lambda {
	return &Lambda{
		Args: args,
		Body: body,
		Env:  e,
	}
}

func (s *Lambda) SexpString() string {
	return "(lambda " + s.Args.SexpString() + " " + s.Body.SexpString() + ")"
}

func (s *Lambda) Equal(o Sexp) bool {
	if oc, ok := o.(*Lambda); ok {
		if !s.Args.Equal(oc.Args) {
			return false
		}

		if !s.Body.Equal(oc.Body) {
			return false
		}

		return true
	}
	return false
}
