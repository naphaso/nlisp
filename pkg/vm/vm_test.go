package vm

import (
	"testing"

	"github.com/naphaso/nlisp/pkg/parser"
	"github.com/stretchr/testify/require"
)

func TestVM(t *testing.T) {
	//vm := New()
	e, err := parser.ParseString("(plus 1 2)")
	require.NoError(t, err)

	vm := New()
	_, err = vm.Eval(e)
	require.NoError(t, err)
}
