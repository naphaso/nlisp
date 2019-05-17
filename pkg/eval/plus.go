package eval

import (
	"errors"
	"fmt"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func plus(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	var v int64
	for args != nil {
		p, ok := args.(*sexp.Pair)
		if !ok {
			return nil, errors.New("invalid argument list")
		}

		r, err := Eval(p.Head, en)
		if err != nil {
			return nil, err
		}

		if vv, ok := r.(*sexp.Int64); ok {
			v += vv.Value
		} else {
			return nil, errors.New("plus currently supports only int64 arguments, not " + fmt.Sprintf("%T", p.Head))
		}

		args = p.Tail
	}

	return sexp.NewInt64(v), nil
}
