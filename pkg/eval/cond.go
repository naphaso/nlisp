package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func cond(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	var condition, ifTrue, ifFalse sexp.Sexp
	var err error
	condition, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	ifTrue, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	ifFalse, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	result, err := Eval(condition, ev)
	if err != nil {
		return nil, err
	}

	if bresult, ok := result.(*sexp.Bool); ok {
		if bresult.Value {
			return Eval(ifTrue, ev)
		} else {
			return Eval(ifFalse, ev)
		}
	}

	return nil, errors.New("non-bool condition")
}

func popArg(args sexp.Sexp) (sexp.Sexp, sexp.Sexp, error) {
	if p, ok := args.(*sexp.Pair); ok {
		return p.Head, p.Tail, nil
	}

	return nil, nil, errors.New("invalid arguments")
}
