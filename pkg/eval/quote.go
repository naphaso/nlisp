package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func quote(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in quote")
	}

	return p.Head, nil
}
