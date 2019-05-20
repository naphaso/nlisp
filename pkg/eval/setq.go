package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func setq(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in set")
	}

	key := p.Head

	p, ok = p.Tail.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in set")
	}

	if key, ok := key.(sexp.Symbol); ok {
		value, err := Eval(p.Head, en)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("setq %s = %s\n", key.Name, value.SexpString())
		en.SetGlobal(key, value)
		return nil, nil
	} else {
		return nil, errors.New("first argument in set is not a symbol")
	}
}
