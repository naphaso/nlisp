package eval

import (
	"errors"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func backquote(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	//fmt.Printf("backquote: %s\n", args.SexpString())
	builder := sexp.NewListBuilder()
	for args != nil {
		p, ok := args.(*sexp.Pair)

		if !ok {
			return nil, errors.New("not a list")
		}

		args = p.Tail

		if fpair, ok := p.Head.(*sexp.Pair); ok {
			if fpsym, ok := fpair.Head.(*sexp.Symbol); ok {
				if fpsym.Name == "comma" {
					if spair, ok := fpair.Tail.(*sexp.Pair); ok {
						curr, err := Eval(spair.Head, en)
						if err != nil {
							return nil, err
						}

						builder.Append(curr)

						continue
					}
				}
			}

			curr, err := backquote(p.Head, en)
			if err != nil {
				return nil, err
			}
			builder.Append(curr)

			continue
		}

		builder.Append(p.Head)

	}

	return builder.Build(), nil
}
