package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/naphaso/nlisp/pkg/lexer"

	"github.com/naphaso/nlisp/pkg/sexp"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	checkSame(t, `(a b c)`)
	checkSame(t, `("asdf" "fdsa")`)
	checkSame(t, `(safd (fd (fdfd 123) df dsf) ("asdf" dd))`)
	checkSame(t, `(a . b)`)
	checkSame(t, `(a nil nil nil)`)

	check(t, `()`, `nil`)
	check(t, `((a b c) . nil)`, `((a b c))`)
	check(t, `(a . nil)`, `(a)`)
	check(t, `  (  a  b c   )`, `(a b c)`)

	checkSeq(t, `()()()`, `nil`, `nil`, `nil`)
	checkSeq(t, `(a)(a)(nil)`, `(a)`, `(a)`, `(nil)`)
}

func checkSame(t *testing.T, expr string) {
	check(t, expr, expr)
}

func checkSeq(t *testing.T, expr string, res ...string) {
	p := New(lexer.New(strings.NewReader(expr)))
	for _, value := range res {
		exp, err := p.Read()
		require.NoError(t, err)
		require.Equal(t, value, sexp.ToString(exp))
	}

	_, err := p.Read()
	require.Equal(t, err, io.EOF)
}

func check(t *testing.T, input, expect string) {
	exp, err := ParseString(input)
	require.NoError(t, err)
	require.Equal(t, expect, sexp.ToString(exp))
}
