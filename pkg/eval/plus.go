package eval

import (
	"errors"
	"fmt"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func plus(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	var v int64
	for args != sexp.Nil {
		p, ok := args.(sexp.Pair)
		if !ok {
			return nil, errors.New("invalid argument list")
		}

		r, err := Eval(p.Head, en)
		if err != nil {
			return nil, err
		}

		if vv, ok := r.(sexp.Int64); ok {
			v += int64(vv)
		} else {
			return nil, errors.New("plus currently supports only int64 arguments, not " + fmt.Sprintf("%T", p.Head))
		}

		args = p.Tail
	}

	return sexp.NewInt64(v), nil
}

func optimizePlus(args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	if args == sexp.Nil {
		return sexp.NewInt64(0), nil
	}

	var err error
	args, err = compileList(args, env)
	if err != nil {
		return nil, err
	}

	// extract 1 arg
	argsP, ok := args.(sexp.Pair)
	if !ok {
		return nil, errors.New("plus: invalid arguments")
	}

	arg1 := argsP.Head

	// if 1 arg
	if argsP.Tail == sexp.Nil {
		return compileSexp(arg1, env)
	}

	// extract 2 arg
	argsP, ok = argsP.Tail.(sexp.Pair)
	if !ok {
		return nil, errors.New("plus: invalid arguments")
	}

	arg2 := argsP.Head

	// if 2 arg
	if argsP.Tail == sexp.Nil {
		if a1, ok := arg1.(sexp.Int64); ok {
			if a2, ok := arg2.(sexp.Int64); ok {
				return sexp.NewInt64(int64(a1) + int64(a2)), nil
			}
		}

		if a1, ok := arg1.(sexp.Int64); ok {
			if a2, ok := arg2.(sexp.Symbol); ok {
				return sexp.NewCompiled(&plusFixed{a: a2, b: int64(a1)}), nil
			}
		}

		if a1, ok := arg1.(sexp.Symbol); ok {
			if a2, ok := arg2.(sexp.Int64); ok {
				return sexp.NewCompiled(&plusFixed{a: a1, b: int64(a2)}), nil
			}
		}

		return sexp.NewCompiled(&plusTwo{a: arg1, b: arg2}), nil
	}

	return sexp.NewPair(sexp.NewSymbol("plus"), args), nil
}

type plusTwo struct {
	a, b sexp.Sexp
}

func (s *plusTwo) Eval(env *sexp.Env) (sexp.Sexp, error) {
	a1, err := Eval(s.a, env)
	if err != nil {
		return nil, err
	}

	a2, err := Eval(s.b, env)
	if err != nil {
		return nil, err
	}

	if a1i, ok := a1.(sexp.Int64); ok {
		if a2i, ok := a2.(sexp.Int64); ok {
			return sexp.NewInt64(int64(a1i) + int64(a2i)), nil
		}
	}

	return nil, errors.New("plusTwo: invalid arguments")
}

type plusFixed struct {
	a sexp.Symbol
	b int64
}

func (s *plusFixed) Eval(env *sexp.Env) (sexp.Sexp, error) {
	if v, ok := env.Get(s.a).(sexp.Int64); ok {
		return sexp.NewInt64(int64(v) + s.b), nil
	}

	return nil, errors.New("plusFixed invalid args")
}

type plusTwoSymbols struct {
	a, b sexp.Symbol
}

func (s *plusTwoSymbols) Eval(env *sexp.Env) (sexp.Sexp, error) {
	v1, ok := env.Get(s.a).(sexp.Int64)
	if !ok {
		return nil, errors.New("plusTwoSymbols invalid args")
	}

	v2, ok := env.Get(s.b).(sexp.Int64)
	if !ok {
		return nil, errors.New("plusTwoSymbols invalid args")
	}

	return sexp.NewInt64(int64(v1) + int64(v2)), nil
}
