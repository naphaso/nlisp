package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func compileLambda(args sexp.Sexp, ev *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in timeit")
	}

	var target *sexp.Lambda
	if lam, ok := p.Head.(*sexp.Lambda); ok {
		target = lam
	} else {
		return nil, errors.New("compile non lambda")
	}

	_ = target
	panic("not implemented")
}

var ErrCannotCompile = errors.New("non-compilable")

func compileCall(f, args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	//fmt.Printf("compile call %s : %s\n", sexp.ToString(f), sexp.ToString(args))
	switch fc := f.(type) {
	case sexp.Symbol:
		// optimize to function call
		fu := env.Get(fc)
		if fu == nil {
			return nil, errors.New("called function does not exist at compile time")
		}

		switch fun := fu.(type) {
		case *sexp.Func:
			if fun.Compile == nil {
				return nil, ErrCannotCompile
			}

			funo, err := fun.Compile(args, env)
			if err != nil {
				return nil, err
			}
			return funo, nil
		case *sexp.Macro:
			e, err := evalMacro(fun, args, env)
			if err != nil {
				return nil, err
			}
			return compileSexp(e, env)
		default:
			return nil, ErrCannotCompile
		}
	default:
		return nil, ErrCannotCompile
	}
}

func compileList(list sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	b := sexp.NewListBuilder()
	for list != sexp.Nil {
		p, ok := list.(sexp.Pair)
		if !ok {
			return nil, errors.New("not a list")
		}

		headOpt, err := compileSexp(p.Head, env)
		if err != nil {
			return nil, err
		}

		b.Append(headOpt)
		list = p.Tail
	}

	return b.Build(), nil
}

func compile(args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	//fmt.Printf("compile : %s\n", sexp.ToString(args))
	var err error
	args, err = EvalList(args, env)
	if err != nil {
		return nil, err
	}

	p, ok := args.(sexp.Pair)
	if !ok {
		return nil, errors.New("compile: invalid arguments")
	}

	if p.Tail != sexp.Nil {
		return nil, errors.New("compile: invalid arguments")
	}

	return compileSexp(p.Head, env)
}

func compileSexp(e sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	//fmt.Printf("compileSexp : %s\n", sexp.ToString(e))
	switch ev := e.(type) {
	case sexp.Pair:
		r, err := compileCall(ev.Head, ev.Tail, en)
		if err == ErrCannotCompile {
			return e, nil
		}
		return r, err
	case *sexp.Lambda:
		o, err := compileSexp(ev.Body, en)
		if err != nil {
			return nil, err
		}
		return sexp.NewLambda(ev.Args, o, ev.Env), nil
	default:
		return e, nil
	}
}
