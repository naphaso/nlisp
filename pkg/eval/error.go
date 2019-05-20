package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func funcError(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	if args == nil {
		return nil, errors.New("unknown error")
	}

	p, ok := args.(sexp.Pair)
	if !ok {
		return nil, errors.New("unknown error arguments")
	}

	errStr, ok := p.Head.(*sexp.String)
	if !ok {
		return nil, errors.New("unknown error arguments")
	}

	return nil, errors.New("error: " + errStr.Value)
}
