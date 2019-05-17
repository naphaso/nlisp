package lexer

import "io"

type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenLPar
	TokenQLPar
	TokenBLPar // `(
	TokenCLPar // ,(
	TokenRPar
	TokenSymbol
	TokenQSymbol // 'symbol
	TokenCSymbol // ,symbol
	TokenBool
	TokenInteger
	TokenString
	TokenDot
	TokenNil
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer interface {
	GetToken() (Token, error)
	PushData(scanner io.RuneScanner)
}
