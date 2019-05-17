package eval

import "github.com/naphaso/nlisp/pkg/sexp"

func list(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	return EvalList(args, en)
}
