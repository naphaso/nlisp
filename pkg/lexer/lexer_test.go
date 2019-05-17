package lexer

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexerSymbols(t *testing.T) {
	requireTokens(t, "(asdf)", []Token{
		{Type: TokenLPar},
		{Type: TokenSymbol, Value: "asdf"},
		{Type: TokenRPar},
	})

	requireTokens(t, "(asdf fdsa)", []Token{
		{Type: TokenLPar},
		{Type: TokenSymbol, Value: "asdf"},
		{Type: TokenSymbol, Value: "fdsa"},
		{Type: TokenRPar},
	})

	requireTokens(t, " (a ( b(c d )f g   )   ", []Token{
		{Type: TokenLPar},
		{Type: TokenSymbol, Value: "a"},
		{Type: TokenLPar},
		{Type: TokenSymbol, Value: "b"},
		{Type: TokenLPar},
		{Type: TokenSymbol, Value: "c"},
		{Type: TokenSymbol, Value: "d"},
		{Type: TokenRPar},
		{Type: TokenSymbol, Value: "f"},
		{Type: TokenSymbol, Value: "g"},
		{Type: TokenRPar},
	})
}

func TestLexerString(t *testing.T) {
	requireTokens(t, `("aaa")`, []Token{
		{Type: TokenLPar},
		{Type: TokenString, Value: "aaa"},
		{Type: TokenRPar},
	})

	requireTokens(t, `("aaa\n")`, []Token{
		{Type: TokenLPar},
		{Type: TokenString, Value: "aaa\n"},
		{Type: TokenRPar},
	})

	requireTokens(t, `( "aaa\r" ( "bbb\t"))`, []Token{
		{Type: TokenLPar},
		{Type: TokenString, Value: "aaa\r"},
		{Type: TokenLPar},
		{Type: TokenString, Value: "bbb\t"},
		{Type: TokenRPar},
		{Type: TokenRPar},
	})
}

func TestLexerBool(t *testing.T) {
	requireTokens(t, `(true false)`, []Token{
		{Type: TokenLPar},
		{Type: TokenBool, Value: "true"},
		{Type: TokenBool, Value: "false"},
		{Type: TokenRPar},
	})
}

func TestLexerInteger(t *testing.T) {
	requireTokens(t, `(1 12 123)`, []Token{
		{Type: TokenLPar},
		{Type: TokenInteger, Value: "1"},
		{Type: TokenInteger, Value: "12"},
		{Type: TokenInteger, Value: "123"},
		{Type: TokenRPar},
	})
}

func requireTokens(t *testing.T, expr string, requiredTokens []Token) {
	lex := New(strings.NewReader(expr))
	var tokens []Token
	for {
		token, err := lex.GetToken()
		if err != nil {
			if err == io.EOF {
				break
			}

			require.NoError(t, err)
		}

		tokens = append(tokens, token)
	}

	require.Equal(t, requiredTokens, tokens)
}
