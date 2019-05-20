package eval

import (
	"fmt"
	"testing"
	"time"

	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/naphaso/nlisp/pkg/sexp"
	"github.com/stretchr/testify/require"
)

func TestEvalQuote(t *testing.T) {
	check(t, `(list 'a)`, `(a)`)
	check(t, `(list '(a b c))'`, `((a b c))`)
}

func TestMacro(t *testing.T) {
	check(t, `
(progn
	(setq inc (macro (v) (list 'setq v (list 'plus v 1))))
	(setq a 0)
	(inc a)
	(inc a)
	(inc a)
	a)
`, `3`)
}

func TestError(t *testing.T) {
	checkErr(t, `(error "hello")`, "error: hello")
}

func TestPrelude(t *testing.T) {
	check(t, `(progn (defun sum (a b) (plus a b)) (sum 100 23))`, `123`)
}

func TestBackquote(t *testing.T) {
	check(t, "(progn (setq a 1) (setq b 2) `(a b))", "(a b)")
	check(t, "(progn (setq a 1) (setq b 2) `(,a ,b))", "(1 2)")
	check(t, "(progn (setq a 1) (setq b 2) `(,(list a b)))", "((1 2))")
	check(t, "(progn (setq a 1) (setq b 2) `((a) (b) ,(list a b) ,(list a) ,(list b) ,a ,b))", "((a) (b) (1 2) (1) (2) 1 2)")
}

// (defn fib (n) cond (< n 2) n (+ (fib (- n 1)) (fib (- n 2))))

//func TestFib(t *testing.T) {
//	check(t, `
//(progn
//	(defun fib (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2)))))
//	(fib 35)
//)
//`, `9227465`)
//}

func TestEvalFunc(t *testing.T) {
	check(t, `(eval plus 2 3))`, `5`)
	check(t, `(eval list 1 2 (plus 1 2))`, `(1 2 3)`)
}

func BenchmarkFib(b *testing.B) {
	env := Global.Wrap()
	//fib, err := parser.ParseString("(defun fib (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2))))))")
	fib, err := parser.ParseString("(setq fib (lambda (n) (cond (less n 2) n (plus (fib (sub n 1)) (fib (sub n 2)))))))")
	require.NoError(b, err)
	_, err = Eval(fib, env)
	require.NoError(b, err)

	call, err := parser.ParseString("(fib 35)")
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t1 := time.Now()
		result, err := Eval(call, env)
		t2 := time.Now()

		fmt.Printf("time: %s\n", t2.Sub(t1))
		require.NoError(b, err)
		require.Equal(b, int64(result.(sexp.Int64)), int64(9227465))
	}
}

func TestEval(t *testing.T) {
	check(t, `(plus)`, `0`)
	check(t, `(plus 1)`, `1`)
	check(t, `(plus 1 2)`, `3`)
	check(t, `(plus 1 2 3)`, `6`)
	check(t, `(plus 1 (plus 2 3))`, `6`)
	check(t, `(plus (plus 1 2) 3)`, `6`)
	check(t, `(plus 1 (plus 2 3))`, `6`)
	check(t, `(plus (plus 1 (plus 2)) 3)`, `6`)

	check(t, `(mul)`, `1`)
	check(t, `(mul 2)`, `2`)
	check(t, `(mul 2 3)`, `6`)
	check(t, `(mul 2 3 4)`, `24`)

	check(t, `(format "%d" 213)`, `"213"`)
	check(t, `(list (setq a 123) (setq b 321) (plus a b))'`, `(nil nil 444)`)
	check(t, `(list (setq a 123) (setq b 321) (format "%d + %d = %d" a b (plus a b)))`, `(nil nil "123 + 321 = 444")`)
	check(t, `((lambda (x) (mul x x)) 2)`, `4`)

	check(t, `(progn (setq a (plus 100 11)) a)`, `111`)

	check(t, `
(progn
	(setq sum (lambda (x)
		(lambda (y)
			(plus x y)
		)
	))
	((sum 123) 321)
)
`, `444`)

	check(t, `
(progn
	(setq sum (lambda (x) (lambda (y) (plus x y))))
	(setq halfsum (sum 123))
	(setq fullsum (halfsum 321)) 
	fullsum)`, `444`)

	check(t, `(progn (setq a 123) (setq b 321) (setq c (plus a b)) c)`, `444`)
	//check(t, `((lambda (list 'x) (list lambda (list 'y) (list mul 'x 'y))) 3) 4)'`, `12`)
	//check(t, `(list (set 'sqr (lambda (x) (list mul x x))) (sqr 2))`, `(nil 4)`)
}

func check(t *testing.T, input string, output string) {
	e, err := parser.ParseString(input)
	require.NoError(t, err)

	r, err := Eval(e, Global.Wrap())
	require.NoError(t, err)

	require.Equal(t, output, r.SexpString())
}

func checkErr(t *testing.T, input string, errStr string) {
	e, err := parser.ParseString(input)
	require.NoError(t, err)

	_, err = Eval(e, Global.Wrap())
	require.Error(t, err)
	require.Equal(t, err.Error(), "error: hello")
}
