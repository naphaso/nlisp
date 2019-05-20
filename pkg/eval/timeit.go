package eval

import (
	"errors"
	"time"

	"github.com/naphaso/nlisp/pkg/sexp"
)

// TODO: implement timeit function to beasure benchmarks
func timeit(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in timeit")
	}

	stmt := p.Head

	startTime := time.Now()
	res, err := Eval(stmt, ev)
	if err != nil {
		return nil, err
	}
	dur := time.Since(startTime)

	return sexp.List(sexp.NewString(dur.String()), res), nil
}
