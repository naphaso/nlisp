package sexp

type Macro struct {
	Args Sexp
	Body Sexp
	Env  *Env
}

func NewMacro(args Sexp, body Sexp, e *Env) *Macro {
	return &Macro{
		Args: args,
		Body: body,
		Env:  e,
	}
}

func (s *Macro) SexpString() string {
	return "(macro " + s.Args.SexpString() + " " + s.Body.SexpString() + ")"
}

func (s *Macro) Equal(o Sexp) bool {
	if oc, ok := o.(*Macro); ok {
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
