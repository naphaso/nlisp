package eval

import (
	"errors"
	"fmt"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func Eval(e sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	//fmt.Printf("eval %s, env:%s\n", e.SexpString(), en.String())
	switch ev := e.(type) {
	case *sexp.Pair:
		op, err := Eval(ev.Head, en)
		if err != nil {
			return nil, err
		}

		switch opc := op.(type) {
		case sexp.Func:
			return opc(ev.Tail, en)
		case *sexp.Lambda:
			return evalLambda(opc, ev.Tail, en)
		case *sexp.Macro:
			e, err := evalMacro(opc, ev.Tail, en)
			if err != nil {
				return nil, err
			}
			return Eval(e, en)
		default:
			return nil, errors.New("is not a function: " + fmt.Sprintf("%T", op))
		}
	case *sexp.Symbol:
		//fmt.Printf("symbol %s = %s\n", ev.Name, sexp.ToString(en.Get(ev)))
		return en.Get(ev), nil
	default:
		return e, nil
	}
}

func evalLambda(lam *sexp.Lambda, args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	lambdaEnv := lam.Env.Wrap()
	dargs := lam.Args
	for args != nil && dargs != nil {
		arg, ok := args.(*sexp.Pair)
		if !ok {
			return nil, errors.New("lambda arguments fail")
		}

		darg, ok := dargs.(*sexp.Pair)
		if !ok {
			return nil, errors.New("lambda arguments fail")
		}

		dargSym, ok := darg.Head.(*sexp.Symbol)
		if !ok {
			return nil, errors.New("lambda signature fail")
		}

		argValue, err := Eval(arg.Head, en)
		if err != nil {
			return nil, err
		}

		lambdaEnv.SetLocal(dargSym, argValue)

		args = arg.Tail
		dargs = darg.Tail
	}

	res, err := Eval(lam.Body, lambdaEnv)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func evalMacro(lam *sexp.Macro, args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	macroEnv := lam.Env.Wrap()
	dargs := lam.Args
	for args != nil && dargs != nil {
		arg, ok := args.(*sexp.Pair)
		if !ok {
			return nil, errors.New("macro arguments fail")
		}

		darg, ok := dargs.(*sexp.Pair)
		if !ok {
			return nil, errors.New("macro arguments fail")
		}

		dargSym, ok := darg.Head.(*sexp.Symbol)
		if !ok {
			return nil, errors.New("macro signature fail")
		}

		argValue := arg.Head

		macroEnv.SetLocal(dargSym, argValue)

		args = arg.Tail
		dargs = darg.Tail
	}

	res, err := Eval(lam.Body, macroEnv)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func EvalList(list sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	arr := make([]sexp.Sexp, 0, 10)
	for list != nil {
		p, ok := list.(*sexp.Pair)
		if !ok {
			return nil, errors.New("not a list")
		}

		v, err := Eval(p.Head, en)
		if err != nil {
			return nil, err
		}

		arr = append(arr, v)
		list = p.Tail
	}

	return sexp.NewList(arr), nil
}
