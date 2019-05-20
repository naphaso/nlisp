package eval

import (
	"errors"
	"fmt"

	"github.com/naphaso/nlisp/pkg/sexp"
)

func format(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in set")
	}

	r, err := Eval(p.Head, en)
	if err != nil {
		return nil, err
	}

	var format string
	if head, ok := r.(*sexp.String); ok {
		format = head.Value
	} else {
		return nil, errors.New("format string is not a string")
	}

	args = p.Tail

	var fargs []interface{}
	for args != nil {
		p, ok = args.(*sexp.Pair)
		if !ok {
			return nil, errors.New("invalid argument list")
		}

		r, err = Eval(p.Head, en)

		switch v := r.(type) {
		case sexp.Int64:
			fargs = append(fargs, int64(v))
		case *sexp.String:
			fargs = append(fargs, v.Value)
		default:
			fargs = append(fargs, v.SexpString())
		}

		args = p.Tail
	}

	return sexp.NewString(fmt.Sprintf(format, fargs...)), nil
}

func funcPrint(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	p, ok := args.(*sexp.Pair)
	if !ok {
		return nil, errors.New("invalid arguments in set")
	}

	r, err := Eval(p.Head, en)
	if err != nil {
		return nil, err
	}

	var data string
	if head, ok := r.(*sexp.String); ok {
		data = head.Value
	} else {
		data = sexp.ToString(r)
	}

	fmt.Print(data)
	return nil, nil
}

func funcPrintln(args sexp.Sexp, en *sexp.Env) (sexp.Sexp, error) {
	_, err := funcPrint(args, en)
	if err != nil {
		return nil, err
	}

	fmt.Println()
	return nil, nil
}
