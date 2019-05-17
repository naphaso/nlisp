package symbol

import "github.com/naphaso/nlisp/pkg/sexp"
import "go.uber.org/atomic"

var registryCodes = map[uint64]*sexp.Symbol{}
var registryNames = map[string]*sexp.Symbol{}
var registryCounter = atomic.Uint64{}

func ByName(name string) *sexp.Symbol {
	if s, ok := registryNames[name]; ok {
		return s
	}

	s := sexp.NewSymbol(registryCounter.Inc(), name)
	registryNames[s.Name] = s
	registryCodes[s.Code] = s
	return s
}

func ByCode(code uint64) *sexp.Symbol {
	return registryCodes[code]
}
