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

func optimizeCond(args sexp.Sexp, env *sexp.Env) (sexp.Sexp, error) {
	var err error
	args, err = compileList(args, env)
	if err != nil {
		return nil, err
	}

	condition, ifTrue, ifFalse, err := sexp.ExtractThreeArgs(args)
	if err != nil {
		return nil, err
	}

	v := &condSimple{
		condition: condition,
		ifTrue:    ifTrue,
		ifFalse:   ifFalse,
	}

	if x, ok := condition.(*sexp.Compiled); ok {
		v.conditionCode = x.Code
	}
	if x, ok := ifTrue.(*sexp.Compiled); ok {
		v.ifTrueCode = x.Code
	}
	if x, ok := ifFalse.(*sexp.Compiled); ok {
		v.ifFalseCode = x.Code
	}

	return sexp.NewCompiled(v), nil
}

type condSimple struct {
	conditionCode sexp.Evaluable
	condition     sexp.Sexp

	ifTrueCode sexp.Evaluable
	ifTrue     sexp.Sexp

	ifFalseCode sexp.Evaluable
	ifFalse     sexp.Sexp
}

func (s *condSimple) Eval(env *sexp.Env) (sexp.Sexp, error) {
	var res sexp.Sexp
	var err error

	if s.conditionCode != nil {
		res, err = s.conditionCode.Eval(env)
	} else {
		res, err = Eval(s.condition, env)
	}
	if err != nil {
		return nil, err
	}

	if resb, ok := res.(*sexp.Bool); ok {
		if resb.Value {
			if s.ifTrueCode != nil {
				return s.ifTrueCode.Eval(env)
			} else {
				return Eval(s.ifTrue, env)
			}
		} else {
			if s.ifFalseCode != nil {
				return s.ifFalseCode.Eval(env)
			} else {
				return Eval(s.ifFalse, env)
			}
		}
	}

	return nil, errors.New("cond non-bool condition")
}

func popArg(args sexp.Sexp) (sexp.Sexp, sexp.Sexp, error) {
	if p, ok := args.(*sexp.Pair); ok {
		return p.Head, p.Tail, nil
	}

	return nil, nil, errors.New("invalid arguments")
}
