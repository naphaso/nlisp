package evalstack

import (
	"testing"

	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/naphaso/nlisp/pkg/sexp"
	"github.com/stretchr/testify/require"
)

func TestEvalStack(t *testing.T) {
	checkS(t, "(plus 2 3)", "5")
	checkS(t, "(plus (plus 1 2) (plus 3 (plus 4 5 6)))", "21")
}

func checkS(t *testing.T, expr string, result string) {
	e, err := parser.ParseString(expr)
	require.NoError(t, err)

	ev := NewEvalStack()
	r, err := ev.Eval(e, sexp.NewEnv())
	require.NoError(t, err)
	require.Equal(t, result, sexp.ToString(r))
}
