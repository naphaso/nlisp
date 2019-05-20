package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func progn(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	var r sexp.Sexp
	var err error
	//fmt.Printf("progn: %s\n", args.SexpString())
	for args != sexp.Nil {
		p, ok := args.(sexp.Pair)
		if !ok {
			return nil, errors.New("invalid argument list")
		}

		//fmt.Printf("progn eval: %s, env:%s\n", p.Head.SexpString(), en.String())
		r, err = Eval(p.Head, en)
		if err != nil {
			return nil, err
		}

		args = p.Tail
	}

	return r, nil
}
