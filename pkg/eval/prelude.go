package eval

import (
	"io"
	"strings"

	"github.com/naphaso/nlisp/pkg/lexer"
	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/naphaso/nlisp/pkg/sexp"
	"github.com/naphaso/nlisp/pkg/symbol"
)

var Global *sexp.Env
var PreludeOld = `
(setq defmacro (macro (name args body) (list 'setq name (list 'macro args body))))
(defmacro defun (name args body) (list 'setq name (list 'lambda args body)))
`

var Prelude = "(setq defmacro (macro (name args body) `(setq ,name (macro ,args ,body)))) (defmacro defun (name args body) `(setq ,name (lambda ,args ,body)))"

func init() {
	Global = sexp.NewEnv()
	reg := func(name string, f func(sexp.Sexp, *sexp.Env) (sexp.Sexp, error)) {
		Global.SetLocal(symbol.ByName(name), sexp.NewFunc(f))
	}

	reg("plus", plus)
	reg("mul", mul)
	reg("list", list)

	reg("setq", setq)
	reg("format", format)
	reg("lambda", lambda)
	reg("progn", progn)

	reg("quote", quote)
	reg("backquote", backquote)
	reg("macro", macro)
	reg("error", funcError)

	reg("cond", cond)
	reg("equal", equal)
	reg("less", less)
	reg("sub", sub)
	reg("print", funcPrint)
	reg("println", funcPrintln)

	p := parser.New(lexer.New(strings.NewReader(Prelude)))
	for {
		e, err := p.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		_, err = Eval(e, Global)
		if err != nil {
			panic(err)
		}
	}
}
