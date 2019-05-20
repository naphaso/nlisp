package sexp

import (
	"go.uber.org/atomic"
)

var registryCodes = map[Symbol]string{}
var registryNames = map[string]Symbol{}
var registryCounter = atomic.Uint64{}

type Symbol uint64

func NewSymbol(name string) Symbol {
	if s, ok := registryNames[name]; ok {
		return Symbol(s)
	}

	s := Symbol(registryCounter.Inc())
	registryNames[name] = s
	registryCodes[s] = name
	return s
}

func (s Symbol) Name() string {
	return registryCodes[s]
}

func (s Symbol) SexpString() string {
	return s.Name()
}

func (s Symbol) Equal(o Sexp) bool {
	if oc, ok := o.(Symbol); ok {
		return oc == s
	}

	return false
}
