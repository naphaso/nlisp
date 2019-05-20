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

	if aInt, ok := a.(sexp.Int64); ok {
		if bInt, ok := b.(sexp.Int64); ok {
			return sexp.NewBool(aInt < bInt), nil
		}
	}

	return nil, errors.New("non-integer compatison")
}

func optimizeLess(args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	var err error
	args, err = compileList(args, env)
	if err != nil {
		return nil, err
	}

	arg1, arg2, err := sexp.ExtractTwoArgs(args)
	if err != nil {
		return nil, err
	}

	if a1, ok := arg1.(sexp.Symbol); ok {
		if a2, ok := arg2.(sexp.Int64); ok {
			return sexp.NewCompiled(&lessSimple{a: a1, b: int64(a2)}), nil
		}
	}

	return sexp.NewPair(sexp.NewSymbol("less"), args), nil
}

type lessSimple struct {
	a sexp.Symbol
	b int64
}

func (s *lessSimple) Eval(env *sexp.Env) (sexp.Sexp, error) {
	if av, ok := env.Get(s.a).(sexp.Int64); ok {
		return sexp.NewBool(int64(av) < s.b), nil
	}

	return nil, errors.New("lessSimple: invalid argumets")
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

	if aInt, ok := a.(sexp.Int64); ok {
		if bInt, ok := b.(sexp.Int64); ok {
			return sexp.NewInt64(int64(aInt) - int64(bInt)), nil
		}
	}

	return nil, errors.New("non-integer sub")
}

func optimizeSub(args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	var err error

	args, err = compileList(args, env)
	if err != nil {
		return nil, err
	}

	arg1, arg2, err := sexp.ExtractTwoArgs(args)
	if err != nil {
		return nil, err
	}

	if a1, ok := arg1.(sexp.Int64); ok {
		if a2, ok := arg2.(sexp.Int64); ok {
			return sexp.NewInt64(int64(a1) - int64(a2)), nil
		}
	}

	if a1, ok := arg1.(sexp.Symbol); ok {
		if a2, ok := arg2.(sexp.Int64); ok {
			return sexp.NewCompiled(&subFixed{a: a1, b: int64(a2)}), nil
		}
	}

	return sexp.NewPair(sexp.NewSymbol("sub"), args), nil
}

type subFixed struct {
	a sexp.Symbol
	b int64
}

func (s *subFixed) Eval(env *sexp.Env) (sexp.Sexp, error) {
	if v, ok := env.Get(s.a).(sexp.Int64); ok {
		return sexp.NewInt64(int64(v) - s.b), nil
	}
	return nil, errors.New("subFixed invalid arguments")
}
