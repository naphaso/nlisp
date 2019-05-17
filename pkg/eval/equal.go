package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func equal(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	var a, b sexp.Sexp
	var err error
	a, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	a, err = Eval(a, ev)
	if err != nil {
		return nil, err
	}

	b, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	b, err = Eval(b, ev)
	if err != nil {
		return nil, err
	}

	return sexp.NewBool(a.Equal(b)), nil
}

func less(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	var a, b sexp.Sexp
	var err error
	a, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	a, err = Eval(a, ev)
	if err != nil {
		return nil, err
	}

	b, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	b, err = Eval(b, ev)
	if err != nil {
		return nil, err
	}

	if aInt, ok := a.(*sexp.Int64); ok {
		if bInt, ok := b.(*sexp.Int64); ok {
			return sexp.NewBool(aInt.Value < bInt.Value), nil
		}
	}

	return nil, errors.New("non-integer compatison")
}

func sub(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	var a, b sexp.Sexp
	var err error
	a, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	a, err = Eval(a, ev)
	if err != nil {
		return nil, err
	}

	b, args, err = popArg(args)
	if err != nil {
		return nil, err
	}

	b, err = Eval(b, ev)
	if err != nil {
		return nil, err
	}

	if aInt, ok := a.(*sexp.Int64); ok {
		if bInt, ok := b.(*sexp.Int64); ok {
			return sexp.NewInt64(aInt.Value - bInt.Value), nil
		}
	}

	return nil, errors.New("non-integer sub")
}
