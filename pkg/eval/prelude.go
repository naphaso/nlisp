package eval

import (
	"io"
	"strings"

	"github.com/naphaso/nlisp/pkg/lexer"
	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/naphaso/nlisp/pkg/sexp"
)

var Global *sexp.Env
var PreludeOld = `
(setq defmacro (macro (name args body) (list 'setq name (list 'macro args body))))
(defmacro defun (name args body) (list 'setq name (list 'lambda args body)))
`

var Prelude = "(setq defmacro (macro (name args body) `(setq ,name (macro ,args ,body)))) (defmacro defun (name args body) `(setq ,name (lambda ,args ,body)))"

func init() {
	Global = sexp.NewEnv()
	reg := func(name string, f sexp.RawFunc, c sexp.CompileFunc) {
		Global.SetLocal(sexp.NewSymbol(name), sexp.NewFunc(f, c))
	}

	reg("plus", plus, optimizePlus)
	reg("mul", mul, nil)
	reg("list", list, nil)

	reg("setq", setq, nil)
	reg("format", format, nil)
	reg("lambda", lambda, nil)
	reg("progn", progn, nil)

	reg("quote", quote, nil)
	reg("backquote", backquote, nil)
	reg("macro", macro, nil)
	reg("error", funcError, nil)

	reg("cond", cond, optimizeCond)
	reg("equal", equal, nil)
	reg("less", less, optimizeLess)
	reg("sub", sub, optimizeSub)
	reg("print", funcPrint, nil)
	reg("println", funcPrintln, nil)

	reg("timeit", timeit, nil)

	reg("compile", compile, nil)

	reg("eval", Eval, nil)

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
