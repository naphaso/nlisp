package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"

	"github.com/naphaso/nlisp/pkg/lexer"

	"github.com/naphaso/nlisp/pkg/eval"
	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/naphaso/nlisp/pkg/sexp"
	"github.com/sirupsen/logrus"
)

const fib1prog = "(defun fib (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2))))))"
const fib2prog = "(setq fib (lambda (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2)))))))"

func runFile(filename string) {
	env := eval.Global.Wrap()
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	p := parser.New(lexer.New(bufio.NewReader(fd)))
	for {
		exp, err := p.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		_, err = eval.Eval(exp, env)
		if err != nil {
			panic(err)
		}
	}
}

type Completer struct {
	Env  *sexp.Env
	Data []prompt.Suggest
}

func (s *Completer) Refresh() {
	data := []prompt.Suggest{
		{
			Text:        "exit",
			Description: "Exit REPL",
		},
		{
			Text:        "quit",
			Description: "Quit REPL",
		},
		{
			Text:        ":q!",
			Description: "Kill this shit",
		},
	}
	for sym, value := range s.Env.GetRawData() {
		var desc string
		switch v := value.(type) {
		case *sexp.Func:
			desc = "function"
		case *sexp.Lambda:
			desc = v.SexpString()
		case *sexp.Macro:
			desc = v.SexpString()
		case sexp.Int64:
			desc = fmt.Sprintf("int64 %d", int64(v))
		case *sexp.String:
			desc = fmt.Sprintf("string %q", v.Value)
		case *sexp.Bool:
			desc = fmt.Sprintf("bool %b", v.Value)
		case *sexp.Pair:
			desc = v.SexpString()
		case *sexp.Float64:
			desc = fmt.Sprintf("float64 %f", v)
		case *sexp.Symbol:
			desc = fmt.Sprintf("symbol %s", v.Name)
		case *sexp.Rune:
			desc = fmt.Sprintf("rune %c", v.Value)
		}

		data = append(data, prompt.Suggest{
			Text:        sym.Name(),
			Description: desc,
		})
	}
	s.Data = data
}

func (s *Completer) Complete(doc prompt.Document) []prompt.Suggest {
	word := doc.GetWordBeforeCursor()
	for strings.HasPrefix(word, "(") {
		word = strings.TrimPrefix(word, "(")
	}
	for strings.HasPrefix(word, "'") {
		word = strings.TrimPrefix(word, "'")
	}
	for strings.HasPrefix(word, ",") {
		word = strings.TrimPrefix(word, ",")
	}

	return prompt.FilterHasPrefix(s.Data, word, true)
}

func main() {
	if len(os.Args) > 1 {
		runFile(os.Args[1])
	} else {
		// TODO: check tty
		runPrompt()
	}
}

func runPrompt() {
	env := eval.Global.Wrap()
	completer := &Completer{Env: env}
	completer.Refresh()

	lex := lexer.New(strings.NewReader(""))
	par := parser.New(lex)

loop:
	for {
		var pt string
		if par.IsEmpty() {
			pt = "nlisp> "
		} else {
			pt = "... "
		}
		line := prompt.Input(pt, completer.Complete)

		if line == ":q!" {
			fmt.Println("Okay ;(")
			return
		}
		if line == "exit" || line == "quit" {
			return
		}

		lex.PushData(strings.NewReader(line + "\n"))
		for {
			e, err := par.Read()
			if err != nil {
				if err == io.EOF {
					continue loop
				}
				logrus.WithError(err).Error("failed to parse expression")
				continue loop
			}

			res, err := eval.Eval(e, env)
			if err != nil {
				logrus.WithError(err).Error("failed to eval")
				continue loop
			}

			if res != nil {
				fmt.Println(res.SexpString())
			}
			completer.Refresh()
		}
	}
}
