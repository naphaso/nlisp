package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func lambda(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in lambda")
	}

	largs := p.Head

	p, ok = p.Tail.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in lambda")
	}

	return sexp.NewLambda(largs, p.Head, en), nil
}
